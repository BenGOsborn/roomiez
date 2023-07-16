package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type SubscriptionRecord struct {
	ID           string        `json:"id"`
	SearchParams *SearchParams `json:"searchParams"`
	Timestamp    time.Time     `json:"timestamp"`
	Email        string        `json:"email"`
}

// Generate key for subscription records
func GenerateKey(email string, secret string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(fmt.Sprint(email, ":", secret))); err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)

	return hex.EncodeToString(hashSum), nil
}
