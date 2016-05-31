package pgwp

import (
	"log"
	"testing"
	"time"

	"gopkg.in/ory-am/dockertest.v2"

	"database/sql"

	_ "github.com/lib/pq"
)

// Item struct.
type Item struct {
	ID   string         `db:"id"`
	Name sql.NullString `db:"name"`
}

func TestConnect(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	// must fail: workers, queue len
	{
		expected := "Could not set up PostgreSQL container."
		var err error
		var x *Pool
		c, err := dockertest.ConnectToPostgreSQL(0, time.Second, func(url string) bool {
			// Check if postgres is responsive...
			x, err = Connect("postgres", url+"bad", 10, 10)
			return err == nil
		})
		if err == nil {
			log.Printf("expected [%s] actual [%s]", expected, err)
			t.Fail()
		}
		c.KillRemove()
	}
	// must fail: bad connection
	{
		expected := "Could not set up PostgreSQL container."
		var err error
		var x *Pool
		c, err := dockertest.ConnectToPostgreSQL(0, time.Second, func(url string) bool {
			// Check if postgres is responsive...
			x, err = Connect("postgres", url, -1, -1)
			return err == nil
		})
		if err == nil {
			log.Printf("expected [%s] actual [%s]", expected, err)
			t.Fail()
		}
		c.KillRemove()
	}
	// Get and Select
	{
		var err error
		var x *Pool
		c, err := dockertest.ConnectToPostgreSQL(4, 5*time.Second, func(url string) bool {
			// Check if postgres is responsive...
			x, err = Connect("postgres", url, 10, 10)
			return err == nil
		})
		if err != nil {
			log.Fatalf("Could not connect to database: %s", err)
			t.Fail()
		}
		prepare(x)

		var one Item
		err = x.Get(&one, "SELECT * FROM mockdata WHERE name=$1", "two")
		if err != nil {
			log.Printf("Get err [%s]", err)
			t.Fail()
		}

		var list []*Item
		err = x.Select(&list, "SELECT * FROM mockdata LIMIT 2")
		if err != nil {
			log.Printf("Select err [%s]", err)
			t.Fail()
		}
		if len(list) != 2 {
			t.Fail()
		}
		log.Printf("Get result [%#v] Select result [%#v]", one, list)
		x.Close()
		c.KillRemove()
	}
}

// prepare generate mock data.
func prepare(pool *Pool) {
	qr := pool.MustExec(`
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
	INSERT INTO mockdata (name) VALUES('three');
	`)
	_, err := qr.LastInsertId()
	if err != nil {
		log.Printf("prepare :  last insert id err [%s]", err)
	}
	_, err = qr.RowsAffected()
	if err != nil {
		log.Printf("prepare : rows affected err [%s]", err)
	}
}
