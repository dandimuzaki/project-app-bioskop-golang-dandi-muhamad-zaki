package entity

import "time"

type Booking struct {
	Model
	UserID      int
	ScreeningID int
	Seats       []Seat
	Status      string
	ExpiredAt   time.Time
}