package entity

type Studio struct {
	Model
	CinemaID int     `json:"cinema_id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
}