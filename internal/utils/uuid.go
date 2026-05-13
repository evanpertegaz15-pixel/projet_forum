package utils

import (
	"github.com/google/uuid"
)

func NewUUID() string {
    u, err := uuid.NewV7()
    if err != nil {
        return uuid.New().String() // fallback v4
    }
    return u.String()
}