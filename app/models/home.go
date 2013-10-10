package models

type Home struct {
	CurrentWeek *Week
	AllWeeks    []*Week
	Stats       Stats
	AppVersion  string
}

type Stats struct {
	Digital, RentStream, Streaming, NWeeks int
}