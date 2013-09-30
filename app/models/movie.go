package models

type Movie struct {
	Id        int
	Week      string
	Title     string
	Imdb      string
	Rank      int
	Services  map[string]bool
	Streaming int // **** These are ints and not bools because we might want to use the number of services in the future ****
	Rental    int
	Purchase  int // This is instant purchase only, not physical DVDs
	DVD       int // We're not actually displaying DVD data on the site (except in data export), but we want to store it
	All       int // this should be streaming + rental + purchase (exclude DVDs)
}
