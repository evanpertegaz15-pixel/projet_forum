package config

import "os"

type Configuration struct {
	Port	string
	Database	string
	SessionSecret	string // pour sécuriser les cookies de session
	Env	string // pour déterminer si on est en dev ou prod
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key) // renvoie "" si la variable d'environnement n'existe pas
	if value != "" {
		return value
	} else {
		return defaultValue
	}
}

func Load() *Configuration {
	return &Configuration{
		Port:	getEnv("PORT", "8080"),
		Database:	getEnv("DB_PATH", "./forum.db"),
		SessionSecret:	getEnv("SESSION_SECRET", "dev-secret"),
		Env:	getEnv("ENV", "development",),
	}
}