package pgwp

import (
	"database/sql"

	"github.com/jimmy-go/jobq"
	"github.com/jmoiron/sqlx"
)

// Pool struct.
type Pool struct {
	DBS        chan *sqlx.DB
	workers    int
	queue      int
	dispatcher *jobq.Dispatcher
}

// Connect handle sqlx.Connect function and receives workers and queue size.
func Connect(driver, url string, workers, queueLen int) (*Pool, error) {
	d, err := jobq.New(workers, queueLen)
	if err != nil {
		return nil, err
	}
	p := &Pool{
		DBS:        make(chan *sqlx.DB, workers),
		workers:    workers,
		queue:      queueLen,
		dispatcher: d,
	}
	for i := 0; i < workers; i++ {
		db, err := sqlx.Connect(driver, url)
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(2)
		select {
		case p.DBS <- db:
		}
	}

	return p, nil
}

// Close closes all connections.
func (p *Pool) Close() {
	for db := range p.DBS {
		db.Close()
		if len(p.DBS) < 1 {
			p.dispatcher.Stop()
			return
		}
	}
}

// Run generic handle call.
func (p *Pool) Run(fn func(db *sqlx.DB) error) error {
	errc := make(chan error, 1)
	task := func() error {
		db := <-p.DBS
		err := fn(db)
		p.DBS <- db
		errc <- err
		return err
	}
	p.dispatcher.Add(task)
	err := <-errc
	return err
}

// Select wrapper for sqlx.Select.
func (p *Pool) Select(results interface{}, query string, args ...interface{}) error {
	errc := make(chan error, 1)
	task := func() error {
		db := <-p.DBS
		err := db.Select(results, query, args...)
		p.DBS <- db
		errc <- err
		return err
	}
	p.dispatcher.Add(task)
	err := <-errc
	return err
}

// Get wrapper for sqlx.Get.
func (p *Pool) Get(result interface{}, query string, args ...interface{}) error {
	errc := make(chan error, 1)
	task := func() error {
		db := <-p.DBS
		err := db.Get(result, query, args...)
		p.DBS <- db
		errc <- err
		return err
	}
	p.dispatcher.Add(task)
	err := <-errc
	return err
}

// MustExec wrapper for sqlx.MustExec.
func (p *Pool) MustExec(query string, args ...interface{}) sql.Result {
	sqlc := make(chan sql.Result, 1)
	task := func() error {
		db := <-p.DBS
		qr := db.MustExec(query, args...)
		_, err := qr.RowsAffected()
		p.DBS <- db
		sqlc <- qr
		return err
	}
	p.dispatcher.Add(task)
	qr := <-sqlc
	return qr
}
