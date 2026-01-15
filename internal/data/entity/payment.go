package entity

type Payment struct {
	Model
	BookingID     int     `json:"booking_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        string  `json:"amount"`
	Status        string  `json:"status"`
	TransactionID *string `json:"transaction_id"`
}