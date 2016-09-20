// Package pgwp contains a wrapper around sqlx for pooling.
//
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Angel Del Castillo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package pgwp

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	// Timeout controls query execution timeout.
	Timeout = time.Duration(4 * time.Second)

	// ErrTimeout is returned when a query exceeds Timeout execution.
	ErrTimeout   = errors.New("pgwp: operation timeout")
	errResultErr = &ResultErr{Error: ErrTimeout}
)

// ConnectFunc type is called when a new connection is required.
type ConnectFunc func() (*sqlx.DB, error)

// ExecFunc defines behavior for functions that return error
type ExecFunc func(db *sqlx.DB) error

// QueryFunc defines behavior for functions that return sql.Result
type QueryFunc func(db *sqlx.DB) sql.Result

// Worker defines database connection.
//
//
type Worker struct {
	db      *sqlx.DB
	x       <-chan time.Time
	errc    chan error
	resultc chan sql.Result
}

// Exec runs a ExecFunc on worker connection.
func (o *Worker) Exec(fn ExecFunc, timeout time.Duration) error {
	// drain err channel
	drainerrors(o.errc)

	go func() {
		o.errc <- fn(o.db)
	}()
	select {
	case err := <-o.errc:
		return err
	case <-time.After(timeout):
		return ErrTimeout
	}
	return ErrTimeout
}

// Query runs a QueryFunc on worker connection.
func (o *Worker) Query(fn QueryFunc, timeout time.Duration) sql.Result {
	// drain result channel
	drainresults(o.resultc)

	go func() {
		o.resultc <- fn(o.db)
	}()
	select {
	case qr := <-o.resultc:
		return qr
	case <-time.After(timeout):
		return errResultErr
	}
	return errResultErr
}

func exec(p *Pool, fn ExecFunc) error {
	w := <-p.workersc
	err := w.Exec(fn, Timeout)
	p.workersc <- w
	return err
}

func drainerrors(c chan error) {
	for i := 0; i < len(c); i++ {
		<-c
	}
}

func drainresults(c chan sql.Result) {
	for i := 0; i < len(c); i++ {
		<-c
	}
}

func execResult(p *Pool, fn QueryFunc) sql.Result {
	w := <-p.workersc
	qr := w.Query(fn, Timeout)
	p.workersc <- w
	return qr
}

// Pool type contains all workers with his respectives connections.
type Pool struct {
	workers  int
	workersc chan *Worker
	queue    int
	connFunc ConnectFunc
	count    int
}

// Connect handle sqlx.Connect function. Receives workers and queue size.
//
// DEPRECATED. Use Open func.
func Connect(driver, url string, workers, queueLen int) (*Pool, error) {
	dbConnectionFunc := func() (*sqlx.DB, error) {
		db, err := sqlx.Connect(driver, url)
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(2)
		return db, nil
	}
	return Open(dbConnectionFunc, workers, queueLen)
}

// Open return a new pool. Receive a database connection function for later reuse.
//
// workers define number of connections open to database.
// queue defines length channel for queries before block. It's useful for control
// over application demand.
func Open(connectionFunc ConnectFunc, workers, queue int) (*Pool, error) {
	// populate workers
	wc := make(chan *Worker, workers)
	for i := 0; i < workers; i++ {
		db, err := connectionFunc()
		if err != nil {
			return nil, err
		}
		w := &Worker{
			db:      db,
			errc:    make(chan error, 1),
			resultc: make(chan sql.Result, 1),
		}
		wc <- w
	}
	p := &Pool{
		workers:  workers,
		workersc: wc,
		queue:    queue,
		connFunc: connectionFunc,
	}
	return p, nil
}

// Close closes all connections.
func (p *Pool) Close() {
	for i := 0; i < p.workers; i++ {
		w := <-p.workersc
		w.db.Close()
	}
}

// Run purpose is to allow database operations if are not yet implemented.
// for general purposes see other funcs before using it.
func (p *Pool) Run(fn func(db *sqlx.DB) error) error {
	return nil
}

// Select wrapper for sqlx.Select.
func (p *Pool) Select(results interface{}, query string, args ...interface{}) error {
	return exec(p, func(db *sqlx.DB) error { return db.Select(results, query, args...) })
}

// Get wrapper for sqlx.Get.
func (p *Pool) Get(result interface{}, query string, args ...interface{}) error {
	return exec(p, func(db *sqlx.DB) error { return db.Get(result, query, args...) })
}

// ResultErr satisfies sql.Result
//
// It's used when a timeout occurs
type ResultErr struct {
	Error error
}

// LastInsertId satisfies sql.Result
func (r *ResultErr) LastInsertId() (int64, error) {
	return 0, r.Error
}

// RowsAffected satisfies sql.Result
func (r *ResultErr) RowsAffected() (int64, error) {
	return 0, r.Error
}

// MustExec wrapper for sqlx.MustExec.
func (p *Pool) MustExec(query string, args ...interface{}) sql.Result {
	return execResult(p, func(db *sqlx.DB) sql.Result { return db.MustExec(query, args...) })
}
