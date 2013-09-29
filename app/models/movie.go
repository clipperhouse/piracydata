package models

type Movie struct {
	Id     int
	Week   string
	Title  string
	Imdb   string
	Rank   int
	Stream bool
	Rent   bool
	Buy    bool
	Any    bool
}
