package models

type Week struct {
	Movies    []Movie
	Name      string
	Streaming int
	Rental    int
	Purchase  int
	DVD       int
	All       int // Actually just streaming + rental + purchase, we ignore DVD availability for now
}
