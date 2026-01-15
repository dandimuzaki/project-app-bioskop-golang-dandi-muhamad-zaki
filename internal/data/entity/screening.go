package entity

import "time"

type Screening struct {
	Model
	StudioID  int       `json:"studio_id"`
	MovieID   int       `json:"movie_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}