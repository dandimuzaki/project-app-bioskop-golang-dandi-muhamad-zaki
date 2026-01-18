package entity

import "time"

type Ticket struct {
	Model
	BookingID int `json:"booking_id"`
	SeatID    int `json:"seat_id"`
	QRToken   int `json:"qr_token"`
	IssuedAt  *time.Time `json:"issued_at,omitempty"`
}