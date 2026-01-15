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
	Date string `json:"date"`
	Movie MovieResponse `json:"movie"`
	Studio StudioResponse `json:"studio"`
	Screenings []ScreeningResponse `json:"screenings"`
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