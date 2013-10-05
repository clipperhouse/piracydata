package models

type Service struct {
	Id        int    `db:"id"`
	MovieId   int    `db:"movie_id"`
	Name      string `db:"name"`
	Available bool   `db:"available"`
}
