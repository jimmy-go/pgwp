package pgwp

import (
	"errors"

	"github.com/jimmy-go/jobq"
	"github.com/jmoiron/sqlx"
	// init driver
	_ "github.com/lib/pq"
)

var (
	errDefault = errors.New("error making query call")
)

// Pool struct.
type Pool struct {
	DBS        chan *sqlx.DB
	workers    int
	queue      int
	dispatcher *jobq.Dispatcher
}

// Connect handle sqlx.Connect function and receives workers and queue size.
func Connect(url string, workers, queueLen int) (*Pool, error) {
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
		db, err := sqlx.Connect("postgres", url)
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

// Filter executes a query and returns a slice of results.
func (p *Pool) Filter(results interface{}, query string) error {
	errc := make(chan error, 1)
	task := func() error {
		db := <-p.DBS
		err := db.Select(results, query)
		p.DBS <- db
		errc <- err
		return err
	}
	p.dispatcher.Add(task)
	err := <-errc
	return err
}

// Locate executes a query and return one row.
func (p *Pool) Locate(result interface{}, query string) error {
	return nil
}

// Close closes all connections.
func (p *Pool) Close() {
	for db := range p.DBS {
		db.Close()
	}
}

// Insert func
func (p *Pool) Insert(query string) error {
	errc := make(chan error, 1)
	task := func() error {
		db := <-p.DBS
		qr := db.MustExec(query)
		_, err := qr.RowsAffected()
		p.DBS <- db
		errc <- err
		return err
	}
	p.dispatcher.Add(task)
	err := <-errc
	return err
}
