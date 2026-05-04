package main

import (
	"forum-dark-jurassic/internal/config"
	"forum-dark-jurassic/internal/database"
	"log"
)

func main() {
	cfg := config.Load()
	db := database.ConnectDB(cfg.Database)
	defer db.Close()
	database.RunMigrations(db)
	log.Println("Base de données initialisée.")
	database.Seed(db)
	log.Println("Base de données remplie avec données par défaut.")
	startServer(db)
}
