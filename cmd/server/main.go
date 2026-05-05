package main

import (
    "html/template"
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

    http.HandleFunc("/test-user", func(w http.ResponseWriter, r *http.Request) {
        user, err := userModel.FindByID(1) // on teste avec l'utilisateur ID=1
        if err != nil || user == nil {
            http.Error(w, "Impossible de récupérer l'utilisateur", http.StatusInternalServerError)
            return
        }

        tmpl := template.Must(template.ParseFiles("./internal/templates/test_user.html"))
        tmpl.Execute(w, user)
    })

    sessionModel := models.NewSessionModel(db)
    authService := services.NewAuthService(userModel, sessionModel)
    authHandler := handlers.NewAuthHandler(authService)
    
    /*categoryModel := models.NewCategoryModel(db)
    topicModel := models.NewTopicModel(db)
    postModel := models.NewPostModel(db)
    categoryService := services.NewCategoryService(categoryModel)
    topicService := services.NewTopicService(topicModel)
    postService := services.NewPostService(postModel)

    http.HandleFunc("/", handlers.Home)
    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            authHandler.ShowRegister(w, r)
        } else if r.Method == http.MethodPost {
            authHandler.Register(w, r)
        }
    })*/
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            authHandler.ShowLogin(w, r)
        } else if r.Method == http.MethodPost {
            authHandler.Login(w, r)
        }
    })
    //http.HandleFunc("/logout", authHandler.Logout)*/
    http.HandleFunc("/profile", authHandler.Profile)

    //http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    
    log.Println("Serveur lancé sur http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}