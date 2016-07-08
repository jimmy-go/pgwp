package pgwp

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	// Timeout controls query execution time.
	Timeout = time.Duration(4 * time.Second)

	errTimeout   = errors.New("pgwp: timeout")
	errResultErr = &ResultErr{Error: errTimeout}
)

// Worker struct
type Worker struct {
	DB *sqlx.DB
	x  <-chan time.Time
}

// Exec func
func (o *Worker) Exec(fn ExecFunc, timeout time.Duration) error {
	return fn(o.DB)
}

// Query func
func (o *Worker) Query(fn QueryFunc, timeout time.Duration) sql.Result {
	return fn(o.DB)
}

// ConnectFunc type is called when a new connection is required.
type ConnectFunc func() (*sqlx.DB, error)

// Pool struct.
type Pool struct {
	workers  int
	workersc chan *Worker
	queue    int
	connFunc ConnectFunc
	count    int
}

// Connect handle sqlx.Connect function and receives workers and queue size.
//
// TODO; add connectFunc as param
func Connect(driver, url string, workers, queueLen int) (*Pool, error) {
	dbfunc := func() (*sqlx.DB, error) {
		db, err := sqlx.Connect(driver, url)
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(2)
		db.SetMaxOpenConns(2)
		return db, nil
	}
	wc := make(chan *Worker, workers)
	for i := 0; i < workers; i++ {
		db, err := dbfunc()
		if err != nil {
			return nil, err
		}
		w := &Worker{
			DB: db,
		}
		wc <- w
	}
	p := &Pool{
		workers:  workers,
		workersc: wc,
		queue:    queueLen,
		connFunc: dbfunc,
	}
	return p, nil
}

// Close closes all connections.
func (p *Pool) Close() {
	for i := 0; i < p.workers; i++ {
		w := <-p.workersc
		w.DB.Close()
	}
}

// Run generic handle call.
func (p *Pool) Run(fn func(db *sqlx.DB) error) error {
	return nil
}

// ExecFunc defines behavior for Select and Get.
type ExecFunc func(db *sqlx.DB) error

func exec(p *Pool, fn ExecFunc) error {
	w := <-p.workersc
	err := w.Exec(fn, Timeout)
	p.workersc <- w
	return err
}

func execResult(p *Pool, fn QueryFunc) sql.Result {
	w := <-p.workersc
	qr := w.Query(fn, Timeout)
	p.workersc <- w
	return qr
}

// Select wrapper for sqlx.Select.
func (p *Pool) Select(results interface{}, query string, args ...interface{}) error {
	return exec(p, func(db *sqlx.DB) error { return db.Select(results, query, args...) })
}

// Get wrapper for sqlx.Get.
func (p *Pool) Get(result interface{}, query string, args ...interface{}) error {
	return exec(p, func(db *sqlx.DB) error { return db.Get(result, query, args...) })
}

// QueryFunc defines behavior for queries.
type QueryFunc func(db *sqlx.DB) sql.Result

// ResultErr satisfies sql.Result
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
