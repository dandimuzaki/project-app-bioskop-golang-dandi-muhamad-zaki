package entity

type Seat struct {
	Model
	StudioID int    `json:"studio_id"`
	SeatCode string `json:"seat_code"`
}