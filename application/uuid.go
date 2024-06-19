package application

import (
	"github.com/google/uuid"
)

// GenerateUUID function that generates and returns a UUID as a string.
func GenerateUUID() string {
	return uuid.New().String()
}
