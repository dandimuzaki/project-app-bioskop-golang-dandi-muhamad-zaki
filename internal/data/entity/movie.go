package entity

import "time"

type Movie struct {
	Model
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