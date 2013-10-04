package models

import (
	"time"
)

type Week struct {
	Movies    []Movie
	Date      time.Time
	Streaming int
	Rental    int
	Purchase  int
	DVD       int
	All       int // Actually just streaming + rental + purchase, we ignore DVD availability for now
}
