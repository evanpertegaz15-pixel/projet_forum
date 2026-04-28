package main

import (
    "log"
    "forum-dark-jurassic/internal/config"
    "forum-dark-jurassic/internal/database"
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
