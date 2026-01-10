package dto

type PaginationQuery struct {
	Page  int
	Limit int
	All   bool
}