package pgwp

import (
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestConnect(t *testing.T) {
	// must fail: workers, queue len
	{
		_, err := Connect("postgres", "host=192.168.2.6 port=5432 dbname=lostsdb user=postgres password=xx123456", -1, -2)
		if err == nil {
			t.Fail()
		}
	}
	// must fail: bad connection
	{
		_, err := Connect("postgres", "host=192.168.2.6 port=5432 dbname=lostsdb user=postgres password=xx", 10, 10)
		if err == nil {
			t.Fail()
		}
	}
	// normal escenario get
	{
		x, err := Connect("postgres", "host=192.168.2.6 port=5432 dbname=lostsdb user=postgres password=xx123456", 10, 10)
		if err != nil {
			log.Printf("err [%s]", err)
			t.Fail()
		}

		var v struct {
			ID string `db:"id"`
		}
		err = x.Get(&v, "SELECT id FROM users LIMIT 1")
		if err != nil {
			log.Printf("err [%s]", err)
			t.Fail()
		}
		x.Close()
	}
	// normal escenario select
	{
		x, err := Connect("postgres", "host=192.168.2.6 port=5432 dbname=lostsdb user=postgres password=xx123456", 10, 10)
		if err != nil {
			log.Printf("err [%s]", err)
			t.Fail()
		}

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
	}
}
