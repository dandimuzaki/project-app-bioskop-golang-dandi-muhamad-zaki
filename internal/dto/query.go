package dto

type PaginationQuery struct {
	Page  int
	Limit int
	All   bool
}

type MovieQuery struct {
	PaginationQuery
	Genre string
}

type ScreeningQuery struct {
	PaginationQuery
	CinemaID int
	MovieID  int
	Date     string
}