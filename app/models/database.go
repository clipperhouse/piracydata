package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type WeekSet []*Week

var CurrentWeek *Week
var Weeks []*Week

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
			//os.Setenv("PGSSLMODE", "require")
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
	dbmap.AddTableWithName(Week{}, "weeks").SetKeys(true, "Id")
	dbmap.AddTableWithName(Movie{}, "movies").SetKeys(true, "Id")
	dbmap.AddTableWithName(Service{}, "services").SetKeys(true, "Id")
	return
}

/*
func LoadCurrentWeek() {
	log.Println("Starting LoadCurrentWeek")

	week := &Week{}

	dbmap := GetDbMap()
	db := dbmap.Db

	var date time.Time
	db.QueryRow("select distinct week from movies order by week desc limit 1").Scan(&date)
	week.Date = date

	var movies []*Movie
	dbmap.Select(&movies, "select * from movies where week = :week", map[string]interface{}{
		"week": week.Date,
	})

	for _, m := range movies {
		m.Summarize()

		var services []Service
		dbmap.Select(&services, "select * from services where movie_id = :movie_id", map[string]interface{}{
			"movie_id": m.Id,
		})

		servicesMap := make(map[string]bool)
		for _, service := range services {
			servicesMap[service.Name] = service.Available
		}

		m.Services = services
	}
	week.Movies = movies

	week.Summarize()

	CurrentWeek = week

	return
}
*/

func LoadAllWeeks() {
	log.Println("Starting LoadAllWeeks")

	var weeks []*Week

	dbmap := GetDbMap()
	db := dbmap.Db

	var dates []time.Time
	rows, err := db.Query("select distinct date from weeks where is_approved = TRUE order by date desc")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			panic(err)
		}
		dates = append(dates, date)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	var movies []*Movie
	dbmap.Select(&movies, "select * from movies")

	for _, m := range movies {
		m.Summarize()

		var services []Service
		dbmap.Select(&services, "select * from services where movie_id = :movie_id", map[string]interface{}{
			"movie_id": m.Id,
		})

		servicesMap := make(map[string]bool)
		for _, service := range services {
			servicesMap[service.Name] = service.Available
		}

		m.Services = services
		m.ServicesMap = servicesMap
	}

	for _, d := range dates {
		week := &Week{}
		week.Date = d
		for _, m := range movies {
			if m.Week.Equal(d) {
				week.Movies = append(week.Movies, m)
			}
		}
		week.Summarize()
		sort.Sort(week)
		weeks = append(weeks, week)
	}
	// add sorting
	sort.Sort(WeekSet(weeks))
	Weeks = weeks
	CurrentWeek = Weeks[len(Weeks)-1]
	return
}

// sorting functions
func (w WeekSet) Len() int {
	return len(w)
}
func (w WeekSet) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}
func (w WeekSet) Less(i, j int) bool {
	return w[i].Date.Before(w[j].Date)
}

func (w Week) Len() int {
	return len(w.Movies)
}
func (w Week) Swap(i, j int) {
	w.Movies[i], w.Movies[j] = w.Movies[j], w.Movies[i]
}
func (w Week) Less(i, j int) bool {
	return w.Movies[i].Rank < w.Movies[j].Rank
}
