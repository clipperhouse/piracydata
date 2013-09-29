package models

type Week struct {
	Movies []Movie
	Name   string
	Stats  WeekStats
}

type WeekStats struct { // Number of movies that are available in each format
	Streaming int
	Rental    int
	Purchase  int
	DVD       int
	All       int // Actually just streaming + rental + purchase, we ignore DVD availability for now
}