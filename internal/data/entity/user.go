package entity

type User struct {
	Model
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password *string `json:"password,omitempty"`
	Role     string  `json:"role"`
}