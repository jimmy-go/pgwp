package pgwp

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Pool contains sqlx.DB pool.
type Pool struct {
	dbc chan *sqlx.DB
}

// Open return a new pool with buffer capacity of size.
func Open(driver, connectionURL string, size int) (*Pool, error) {
	p := &Pool{
		dbc: make(chan *sqlx.DB, size),
	}
	for i := 0; i < size; i++ {
		db, err := sqlx.Connect(driver, connectionURL)
		if err != nil {
			return nil, err
		}
		p.dbc <- db
	}
	return p, nil
}

// Close closes all connections. Must be called at termination time.
func (p *Pool) Close() error {
	for {
		select {
		case db := <-p.dbc:
			if err := db.Close(); err != nil {
				return err
			}
		}
	}
}

// Execute allows database operations that are not yet implemented.
func (p *Pool) Execute(fn func(db *sqlx.DB) error) error {
	for {
		select {
		case db := <-p.dbc:
			err := fn(db)
			p.dbc <- db
			return err
		}
	}
}

// SelectContext wrapper for sqlx.SelectContext.
func (p *Pool) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.Execute(func(db *sqlx.DB) error {
		return db.SelectContext(ctx, dest, query, args...)
	})
}

// GetContext wrapper for sqlx.GetContext.
func (p *Pool) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.Execute(func(db *sqlx.DB) error {
		return db.GetContext(ctx, dest, query, args...)
	})
}

// ExecContext wrapper for sqlx.ExecContext.
func (p *Pool) ExecContext(ctx context.Context, query string, args ...interface{}) error {
	return p.Execute(func(db *sqlx.DB) error {
		_, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		return nil
	})
}
