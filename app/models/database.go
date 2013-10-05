package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

func openDb() *sql.DB {
	connection := os.Getenv("DATABASE_URL")
	sslmode := os.Getenv("PGSSLMODE") // based on https://github.com/lib/pq/commit/8875df52e9844f4c3fce993c8598bbd1c95c8a0f
	log.Println("sslmode: " + sslmode)

	if sslmode == "" {
		is_heroku := false
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "HEROKU_POSTGRESQL") {
				is_heroku = true
				break
			}
		}
		if is_heroku {
			os.Setenv("PGSSLMODE", "required")
		} else {
			os.Setenv("PGSSLMODE", "disable")
		}
	}
	log.Println(connection)

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
	}

	return db
}

func GetDbMap() (dbmap *gorp.DbMap) {
	db := openDb()
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Movie{}, "movies").SetKeys(true, "Id")
	return
}
