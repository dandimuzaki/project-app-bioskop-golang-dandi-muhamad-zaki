package entity

type Screening struct {
	Model
	StudioID  int    `json:"studio_id"`
	MovieID   int    `json:"movie_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}