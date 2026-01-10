package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

// Generate random token
func GenerateRandomToken(length int) (uuid.UUID, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return uuid.Nil, err
	}
	token := uuid.MustParse(hex.EncodeToString(bytes))
	return token, nil
}
