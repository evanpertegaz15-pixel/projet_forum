package main

import (
    "database/sql"
    "log"
    "net/http"
    "forum-dark-jurassic/internal/config"
    "forum-dark-jurassic/internal/database"
    "forum-dark-jurassic/internal/handlers"
    "forum-dark-jurassic/internal/models"
    "forum-dark-jurassic/internal/services"
)

func main() {
    cfg := config.Load()
    db := database.ConnectDB(cfg.Database)
	defer db.Close()
    database.RunMigrations(db)
    log.Println("Base de données initialisée.")
    database.Seed(db)
    log.Println("Base de données remplie avec données par défaut.")
    
    userModel := models.NewUserModel(db)
    sessionModel := models.NewSessionModel(db)
    authService := services.NewAuthService(userModel, sessionModel)
    authHandler := handlers.NewAuthHandler(authService)

    http.HandleFunc("/register", authHandler.Register)
    http.HandleFunc("/login", authHandler.Login)
    http.HandleFunc("/logout", authHandler.Logout)
    http.HandleFunc("/profile", authHandler.Profile)

    http.Handle("/", http.FileServer(http.Dir("./internal/templates")))

    log.Println("Serveur lancé sur http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
