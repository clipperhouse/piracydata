package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"time"
)

var CurrentWeek *Week

func openDb() *sql.DB {
	connection := os.Getenv("DATABASE_URL")
	sslmode := os.Getenv("PGSSLMODE") // based on https://github.com/lib/pq/commit/8875df52e9844f4c3fce993c8598bbd1c95c8a0f

	if sslmode == "" {
		is_heroku := false
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "HEROKU_POSTGRESQL") {
				is_heroku = true
				break
			}
		}
		if is_heroku {
			os.Setenv("PGSSLMODE", "require")
		} else {
			os.Setenv("PGSSLMODE", "disable")
		}
	}

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
	dbmap.AddTableWithName(Service{}, "services").SetKeys(true, "Id")
	return
}

func LoadCurrentWeek() {
	log.Println("Starting LoadCurrentWeek")

	week := &Week{}

	dbmap := GetDbMap()
	db := dbmap.Db

	var date time.Time
	db.QueryRow("select distinct week from movies order by week desc limit 1").Scan(&date)
	week.Date = date

	var movies []Movie
	dbmap.Select(&movies, "select * from movies where week = :week", map[string]interface{}{
		"week": week.Date,
	})

	for i, _ := range movies {
		movies[i].Summarize()

		var services []Service
		dbmap.Select(&services, "select * from services where movie_id = :movie_id", map[string]interface{}{
			"movie_id": movies[i].Id,
		})

		servicesMap := make(map[string]bool)
		for _, service := range services {
			servicesMap[service.Name] = service.Available
		}

		movies[i].Services = services
	}
	week.Movies = movies

	week.Summarize()

	CurrentWeek = week

	return
}
