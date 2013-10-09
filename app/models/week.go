package models

import (
	"time"
)

type Week struct {
	Movies    []*Movie
	Date      time.Time
	Streaming int
	Rental    int
	Purchase  int
	DVD       int
	All       int // Actually just streaming + rental + purchase, we ignore DVD availability for now
}

func (week *Week) Summarize() {
	var streaming, rental, purchase, dvd, all int
	for m := range week.Movies {
		if week.Movies[m].Streaming > 0 {
			streaming += 1
		}
		if week.Movies[m].Rental > 0 {
			rental += 1
		}
		if week.Movies[m].Purchase > 0 {
			purchase += 1
		}
		if week.Movies[m].DVD > 0 {
			dvd += 1
		}
		if week.Movies[m].All > 0 {
			all += 1
		}
	}
	week.Streaming = streaming
	week.Rental = rental
	week.Purchase = purchase
	week.DVD = dvd
	week.All = all
}
