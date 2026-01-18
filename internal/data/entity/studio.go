package entity

type StudioType struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Row    int     `json:"row"`
	Column int     `json:"column"`
	Price  float64 `json:"price"`
}

type Studio struct {
	Model
	CinemaID int     `json:"cinema_id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Row      *int    `json:"row,omitempty"`
	Column   *int    `json:"column,omitempty"`
}