package entity

type PaymentMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Payment struct {
	Model
	BookingID     int `json:"booking_id"`
	PaymentMethod PaymentMethod
	Amount        string  `json:"amount"`
	Status        string  `json:"status"`
	TransactionID *string `json:"transaction_id"`
}