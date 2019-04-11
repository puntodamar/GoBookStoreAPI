package driver

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

var db *sql.DB

func LogFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func ConnectDB() *sql.DB {
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	LogFatal(err)
	db, err = sql.Open("postgres", pgURL)

	LogFatal(err)
	db.Ping()
	LogFatal(err)
	log.Println(pgURL)

	return db
}
