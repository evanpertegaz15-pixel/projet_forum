package utils
// pour les session.id, les notifications, les logs, les uploads d'images
// NewUUID() string

import (
	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}