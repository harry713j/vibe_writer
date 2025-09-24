package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(size int) (string, error) {
	bytes := make([]byte, size/2)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
