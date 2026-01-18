package dto

import "time"

type CinemaResponse struct {
	CinemaID int `json:"cinema_id"`
	Name string `json:"name"`
	Location string `json:"location"`
}

type StudioResponse struct {
	StudioID int `json:"studio_id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Price float64 `json:"price"`
}

type MovieResponse struct {
	MovieID int `json:"movie_id"`
	Title       string    `json:"title"`
	Synopsis    string    `json:"synopsis"`
	Genres      []string  `json:"genres,omitempty"`
	PosterURL   string    `json:"poster_url"`
	TrailerURL  string    `json:"trailer_url"`
	Duration    string    `json:"duration_minute"`
	ReleaseDate time.Time `json:"release_date"`
	Language    string    `json:"language"`
	RatingAge   string    `json:"rating_age"`
}

type ScreeningResponse struct {
	ScreeningID int `json:"screening_id"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
}

// Get schedule based on selected cinema and date
type MovieScreening struct {
	Movie MovieResponse `json:"movie"`
	Studio StudioResponse `json:"studio"`
	Screenings []ScreeningResponse `json:"screenings"`
}

type MovieByCinema struct {
	Date string `json:"date"`
	List []MovieScreening `json:"list"`
}

type MovieScreeningRow struct {
	Date string `json:"date"`
	Movie MovieResponse `json:"movie"`
	Studio StudioResponse `json:"studio"`
	Screening ScreeningResponse `json:"screening"`
}

type SeatResponse struct {
	ID int `json:"seat_id"`
	SeatCode string `json:"seat_code"`
	Status *string `json:"status,omitempty"`
}

type BookingResponse struct {
	BookingID int `json:"booking_id"`
	BookingDate string `json:"booking_date"`
	Movie MovieResponse `json:"movie"`
	Cinema CinemaResponse `json:"cinema"`
	Studio StudioResponse	`json:"studio"`
	Seats       []SeatResponse `json:"seats"`
	TotalAmount float64 `json:"total_amount"`
	Status      string `json:"status"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type PaymentResponse struct {
	PaymentID int `json:"payment_id"`
	TransactionID *string `json:"transaction_id"`
}

type ProfileResponse struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

type Ticket struct {
	SeatCode string `json:"seat_code"`
	QRToken string `json:"qr_token"`
}

type TicketEmail struct {
	Profile ProfileResponse `json:"profile"`
	BookingID int `json:"booking_id"`
	BookingDate string `json:"booking_date"`
	Movie MovieResponse `json:"movie"`
	Cinema CinemaResponse `json:"cinema"`
	Studio StudioResponse	`json:"studio"`
	Screening ScreeningResponse `json:"screening"`
	Tickets []Ticket `json:"tickets"`
}

type BookingHistory struct {
	ID int `json:"booking_id"`
	MovieTitle string `json:"movie_title"`
	Cinema string `json:"cinema"`
	PosterURL string `json:"poster_url"`
	Seats       []string `json:"seats"`
	Status      string `json:"status"`
	Date string `json:"date"`
	StartTime string `json:"start_time"`
}