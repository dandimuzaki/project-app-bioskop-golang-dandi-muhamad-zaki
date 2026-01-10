package entity

type Cinema struct {
	Model
	Name     string `json:"name"`
	Location string `json:"location"`
}