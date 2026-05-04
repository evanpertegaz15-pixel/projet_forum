// table clé / valeur -> get, set un setting, cache mémoire, reload

package services

import (
	"errors"
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

// Charger tous les settings en mémoire
func (s *SettingService) Load() error {
	settings, err := s.Settings.GetAll()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, setting := range settings {
		s.cache[setting.Key] = setting.Value
	}

	return nil
}

// Get avec cache
func (s *SettingService) Get(key string) (string, error) {
	s.mu.RLock()
	value, ok := s.cache[key]
	s.mu.RUnlock()

	if ok {
		return value, nil
	}

	// fallback DB
	setting, err := s.Settings.GetByKey(key)
	if err != nil {
		return "", err
	}
	if setting == nil {
		return "", errors.New("setting introuvable")
	}

	// update cache
	s.mu.Lock()
	s.cache[key] = setting.Value
	s.mu.Unlock()

	return setting.Value, nil
}

// Set (update ou create)
func (s *SettingService) Set(key, value string) error {

	err := s.Settings.Upsert(key, value)
	if err != nil {
		return err
	}

	// update cache
	s.mu.Lock()
	s.cache[key] = value
	s.mu.Unlock()

	return nil
}

// Reload complet (vider + recharger)
func (s *SettingService) Reload() error {
	s.mu.Lock()
	s.cache = make(map[string]string)
	s.mu.Unlock()

	return s.Load()
}

// Supprimer un setting
func (s *SettingService) Delete(key string) error {
	err := s.Settings.Delete(key)
	if err != nil {
		return err
	}

	s.mu.Lock()
	delete(s.cache, key)
	s.mu.Unlock()

	return nil
}

// Tout récupérer (depuis cache)
func (s *SettingService) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// copie pour éviter modification externe
	copyMap := make(map[string]string)
	for k, v := range s.cache {
		copyMap[k] = v
	}
	return copyMap
}
