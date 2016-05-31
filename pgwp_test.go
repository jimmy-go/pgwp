package pgwp

import (
	"log"
	"testing"
	"time"

	"gopkg.in/ory-am/dockertest.v2"

	_ "github.com/lib/pq"
)

func TestConnect(t *testing.T) {
	// must fail: workers, queue len
	{
		var err error
		var x *Pool
		c, err := dockertest.ConnectToPostgreSQL(2, time.Second, func(url string) bool {
			// Check if postgres is responsive...
			x, err = Connect("postgres", url+"bad", 10, 10)
			return err == nil
		})
		if err != nil {
			log.Fatalf("Could not connect to database: %s", err)
		}
		c.KillRemove()
	}
	// must fail: bad connection
	{
		var err error
		var x *Pool
		c, err := dockertest.ConnectToPostgreSQL(2, time.Second, func(url string) bool {
			// Check if postgres is responsive...
			x, err = Connect("postgres", url, -1, -1)
			return err == nil
		})
		if err == nil {
			t.Fail()
		}
		c.KillRemove()
	}
	// normal escenario get
	{
		var err error
		var x *Pool
		c, err := dockertest.ConnectToPostgreSQL(2, time.Second, func(url string) bool {
			// Check if postgres is responsive...
			x, err = Connect("postgres", url, 10, 10)
			return err == nil
		})
		if err != nil {
			log.Fatalf("Could not connect to database: %s", err)
			t.Fail()
		}
		log.Printf("c [%v]", c)

		var v struct {
			ID string `db:"id"`
		}
		err = x.Get(&v, "SELECT id FROM users LIMIT 1")
		if err != nil {
			log.Printf("err [%s]", err)
			t.Fail()
		}
		x.Close()
		c.KillRemove()
	}
	// normal escenario select
	{
		var err error
		var x *Pool
		c, err := dockertest.ConnectToPostgreSQL(2, time.Second, func(url string) bool {
			// Check if postgres is responsive...
			x, err = Connect("postgres", url, 10, 10)
			return err == nil
		})
		if err != nil {
			log.Fatalf("Could not connect to database: %s", err)
			t.Fail()
		}
		log.Printf("c [%v]", c)

		var v []struct {
			ID string `db:"id"`
		}
		err = x.Select(&v, "SELECT id FROM users LIMIT 2")
		if err != nil {
			log.Printf("err [%s]", err)
			t.Fail()
		}
		if len(v) != 2 {
			t.Fail()
		}
		x.Close()
		c.KillRemove()
	}
}
