package utils

import (
	"fmt"

	"github.com/project-app-bioskop-golang/internal/dto"
)

func SendOTP(data dto.OTPResponse) string {
	return fmt.Sprintf(`
	<h2>Email Verification</h2>

	<p>Hello %s!</p>

	<p>
	Thank you for registering. Please use the verification code below to
	confirm your email address:
	</p>

	<div style='
		font-size: 24px;
		font-weight: bold;
		letter-spacing: 4px;
		margin: 16px 0;
	'>
		%v
	</div>

	<p>
	This code will expire in <strong>5 minutes</strong>.
	</p>

	<p>
	If you did not request this, please ignore this email.
	</p>

	<p style='color: #888; font-size: 12px;'>
	Do not share this code with anyone.
	</p>
	`, data.Name, data.OTP)
}