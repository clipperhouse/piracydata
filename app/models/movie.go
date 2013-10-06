package models

import (
	"time"
)

type Movie struct {
	Id          int             `db:"id"`
	Week        time.Time       `db:"week"`
	Title       string          `db:"title"`
	Imdb        string          `db:"imdb"`
	Rank        int             `db:"rank"`
	Streaming   int             `db:"streaming"`
	Rental      int             `db:"rental"`
	Purchase    int             `db:"purchase"`
	DVD         int             `db:"dvd"`
	All         int             `db:"-"` // ignore, calculate at runtime
	Services    []Service       `db:"-"`
	ServicesMap map[string]bool `db:"-"`
}

func (movie *Movie) Summarize() {
	movie.All = movie.Streaming + movie.Rental + movie.Purchase
}
