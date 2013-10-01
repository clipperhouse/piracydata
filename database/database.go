package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// for local connection, set up postgres with user piracydata, create database piracydata, then:
// export DATABASE_URL="dbname=piracydata host=localhost port=5432 sslmode=disable user=piracydata"
// Heroku uses the DATABASE_URL enivronment variable, though theirs is an actual URL
// can't use URL-style locally with postgres.app, as there is no way to disable SSL

func Open() *sql.DB {
	connection := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
	}
	return db
}
