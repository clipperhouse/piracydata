package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func openDb() *sql.DB {
	connection := os.Getenv("DATABASE_URL")
	sslmode := os.Getenv("PGSSLMODE") // based on https://github.com/lib/pq/commit/8875df52e9844f4c3fce993c8598bbd1c95c8a0f
	if sslmode == "" {
		os.Setenv("PGSSLMODE", "disable")
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
