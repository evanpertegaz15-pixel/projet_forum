package services

import (
	"sync"
	"forum-dark-jurassic/internal/models"
)

type SettingService struct {
	Settings *models.SettingModel

	cache map[string]string
	mu    sync.RWMutex
}

func NewSettingService(settings *models.SettingModel) *SettingService {
	return &SettingService{
		Settings: settings,
		cache:    make(map[string]string),
	}
}