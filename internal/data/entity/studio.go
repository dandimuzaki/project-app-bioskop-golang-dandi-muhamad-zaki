package entity

type Studio struct {
	Model
	CinemaID int    `json:"cinema_id"`
	Name     string `json:"name"`
}