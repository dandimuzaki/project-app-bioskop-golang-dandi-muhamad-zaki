package entity

import "time"

type OTP struct {
	Model
	Email     string `json:"email"`
	OTPHash   string `json:"otp_hash"`
	ExpiredAt time.Time `json:"expired_at"`
	UsedAt time.Time `json:"used_at"`
}