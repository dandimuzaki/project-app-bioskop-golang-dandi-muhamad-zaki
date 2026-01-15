package utils

import (
	"crypto/rand"
	"encoding/base64"
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