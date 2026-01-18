package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

func GenerateTransactionID(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	return encoded[:length], nil
}

// GenerateOTP generates a cryptographically secure numeric OTP of a given length.
func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		b[i] = digits[num.Int64()]
	}
	return string(b), nil
}