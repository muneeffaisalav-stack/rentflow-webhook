package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// GenerateSecureToken generates a secure random token

func GenerateSecureToken(length int) (string, error) {
	log.Printf("Generating secure token of length: %d", length)
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generating secure token: %v", err)
		return "", err
	}
	token := hex.EncodeToString(b)
	log.Printf("Successfully generated secure token")
	return token, nil
}
