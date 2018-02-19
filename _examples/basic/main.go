package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"database/sql"

	"github.com/goinggo/work"
	"github.com/jimmy-go/pgwp"
	"github.com/lib/pq"
)

var (
	pgpool     *pgwp.Pool
	dbHost     = flag.String("host", "", "PosgreSQL host.")
	dbDatabase = flag.String("database", "", "PosgreSQL database.")
	dbUsername = flag.String("u", "", "PosgreSQL username.")
	dbPassword = flag.String("p", "", "PosgreSQL password.")
)

// U struct
type U struct {
	ID        sql.NullString `db:"id"`
	Name      string         `db:"name"`
	Email     string         `db:"email"`
	CreatedAt pq.NullTime    `db:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at"`
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()
	go goroutines()

	s := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", *dbHost, *dbDatabase, *dbUsername, *dbPassword)
	var err error
	pgpool, err = pgwp.Connect("postgres", s, 50, 50)
	if err != nil {
		log.Fatalf(" err [%s]", err)
	}

	pool, err := work.New(600, time.Duration(15*time.Minute), func(string) {})
	if err != nil {
		log.Printf("err [%s]", err)
	}
	for i := 0; i < 1000*1000; i++ {
		pool.Run(&W{ID: i})
		pool.Run(&I{ID: i})
	}
	pool.Shutdown()

	pgpool.Close()

	log.Printf("done T [%s]", time.Since(start))
}

// W safisfies work.Worker
type W struct {
	ID int
}

// Work func.
func (w *W) Work(id int) {
	now := time.Now()
	var users []*U
	err := pgpool.Select(&users, `SELECT * FROM users LIMIT $1`, 100)
	if err != nil {
		log.Printf("Work : err [%s]", err)
	}
	log.Printf("Work id [%v] results [%v] T[%s]", w.ID, len(users), time.Since(now))
}

// I safisfies work.Worker
type I struct {
	ID int
}

// Work func.
func (w *I) Work(id int) {
	err := pgpool.MustExec(`INSERT INTO users (email, name) VALUES ($1, $2)`,
		fmt.Sprintf("email-%v", w.ID), fmt.Sprintf("stress-%v", w.ID))
	if err != nil {
		log.Printf("Work : err [%s]", err)
	}
}

func goroutines() {
	for {
		select {
		case <-time.After(15 * time.Second):
			log.Printf("GOROUTINES [%v]", runtime.NumGoroutine())
		}
	}
}
