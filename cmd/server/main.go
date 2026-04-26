package main

import (
    "log"
    "forum-dark-jurassic/internal/database"
)

func main() {
    db := database.ConnectDB("forum.db")
	defer db.Close()
    database.RunMigrations(db)
    log.Println("Base de données initialisée.")
	//startServer(db) // envoyer la db aux handlers
}
