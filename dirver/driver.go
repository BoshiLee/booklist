package dirver

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ConnectToDB() *sql.DB {
	pqUrl, err := pq.ParseURL(os.Getenv("SQL_URL"))
	logFatal(err)
	db, err = sql.Open("postgres", pqUrl)
	logFatal(err)
	err = db.Ping()
	return db
}
