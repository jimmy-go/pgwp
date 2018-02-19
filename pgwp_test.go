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

func TestConnectAndSelect(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	db, err := connectDB()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = migrateSchema(db)
	assert.Nil(t, err)

	var list []struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}
	ctx := context.TODO()
	err = db.SelectContext(ctx, &list, `SELECT id,name FROM mockdata LIMIT 3`)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(list))
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
