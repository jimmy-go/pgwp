package pgwp

import (
	"context"
	"fmt"
	"log"
	"testing"

	// import driver.
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	dockertest "gopkg.in/ory/dockertest.v3"
)

func connectDB() (*Pool, error) {
	var err error
	var dpool *dockertest.Pool
	dpool, err = dockertest.NewPool("")
	if err != nil {
		return nil, err
	}
	resource, err := dpool.Run("postgres", "9.6", nil)
	if err != nil {
		return nil, err
	}
	var x *Pool
	if err = dpool.Retry(func() error {
		x, err = Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp")), 5)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return x, nil
}

func TestConnect(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	pool, err := connectDB()
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	err = migrateSchema(pool)
	assert.Nil(t, err)
}

// migrateSchema generate mock data.
func migrateSchema(pool *Pool) error {
	ctx := context.TODO()
	err := pool.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE mockdata
		(
			id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
			name text,
			CONSTRAINT mockdata_pkey PRIMARY KEY (id)
		)
		WITH (
			OIDS=FALSE
		);

		INSERT INTO mockdata (name) VALUES('one');
		INSERT INTO mockdata (name) VALUES('two');
		INSERT INTO mockdata (name) VALUES('three') RETURNING id;
	`)
	return err
}
