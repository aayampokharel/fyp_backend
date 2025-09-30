package common

import "github.com/google/uuid"

func GenerateUUID(length int) string {
	return uuid.New().String()[:length]
}
