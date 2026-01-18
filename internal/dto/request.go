package dto

type CinemaRequest struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type StudioRequest struct {
	CinemaID int    `json:"cinema_id" validate:"required,gt=0"`
	Name     string `json:"name" validate:"required"`
	Type     int    `json:"type" validate:"required"`
}

type MovieRequest struct {
	Title       string `json:"title" validate:"required"`
	Synopsis    string `json:"synopsis" validate:"required"`
	Genres      []int  `json:"genres" validate:"required"`
	PosterURL   string `json:"poster_url" validate:"required"`
	TrailerURL  string `json:"trailer_url" validate:"required"`
	Duration    int    `json:"duration_minute" validate:"required,gt=0"`
	ReleaseDate string `json:"release_date" validate:"required,datetime=2006-01-02"`
	Language    string `json:"language" validate:"required"`
	RatingAge   string `json:"rating_age" validate:"required"`
}

type ScreeningRequest struct {
	StudioID   int      `json:"studio_id"`
	MovieID    int      `json:"movie_id"`
	StartDate  string   `json:"start_date" validate:"required,datetime=02-01-2006"`
	EndDate    string   `json:"end_date" validate:"required,datetime=02-01-2006"`
	StartHours []string `json:"start_hours" validate:"required,dive,datetime=15.04"`
}

type UpdateScreeningRequest struct {
	StudioID  int    `json:"studio_id"`
	MovieID   int    `json:"movie_id"`
	StartTime string `json:"start_time" validate:"required"`
}

type BookingRequest struct {
	UserID      int   `json:"user_id" validate:"required"`
	ScreeningID int   `json:"screening_id" validate:"required"`
	Seats       []int `json:"seats" validate:"required"`
}

type PaymentRequest struct {
	BookingID     int     `json:"booking_id" validate:"required"`
	PaymentMethod int     `json:"payment_method" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
}

type UpdatePayment struct {
	PaymentID     int     `json:"payment_id"`
	Status        string  `json:"status"`
	TransactionID *string `json:"transaction_id"`
}

type OTP struct {
	Email   string `json:"email"`
	OTPHash string `json:"otp_hash"`
}

type Attachment struct {
	FileName    string
	FileByte    []byte
	ContentType string
}

type EmailRequest struct {
	From        string
	To          string
	Subject     string
	Body        string
	Attachments []Attachment
}

type StudioType struct {
	Name   string  `json:"name" validate:"required"`
	Row    int     `json:"row" validate:"required,gt=0"`
	Column int     `json:"column" validate:"required,gt=0"`
	Price  float64 `json:"price" validate:"required,gt=0"`
}