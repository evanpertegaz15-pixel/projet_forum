package main

import (
	"forum-dark-jurassic/internal/config"
	"forum-dark-jurassic/internal/database"
	"forum-dark-jurassic/internal/handlers"
	"forum-dark-jurassic/internal/middleware"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/services"
	"forum-dark-jurassic/internal/utils"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	db := database.ConnectDB(cfg.Database)
	defer db.Close()
	database.RunMigrations(db)
	log.Println("Base de données initialisée.")
	database.Seed(db)
	log.Println("Base de données remplie avec données par défaut.")

	err := utils.MinifyCSSFile("static/styles.css", "static/styles.min.css")
	if err != nil {
		log.Println("Erreur minification CSS:", err)
	} else {
		log.Println("CSS minifié avec succès.")
	}

	userModel := models.NewUserModel(db)
	sessionModel := models.NewSessionModel(db)
	postModel := models.NewPostModel(db)
	categoryModel := models.NewCategoryModel(db)
	topicModel := models.NewTopicModel(db)

	userService := services.NewUserService(userModel)
	authService := services.NewAuthService(userModel, sessionModel)
	postService := services.NewPostService(postModel)
	categoryService := services.NewCategoryService(categoryModel)
	topicService := services.NewTopicService(topicModel)

	homeHandler := handlers.NewHomeHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	topicHandler := handlers.NewTopicHandler(topicService, postService, categoryService, authService)
	postHandler := handlers.NewPostHandler(postService, authService)

	// ─── Routes ───────────────────────────────────────────────────────────────
	http.HandleFunc("/", homeHandler.ShowHome)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.ShowRegister(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.Register(w, r)
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.ShowLogin(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.Login(w, r)
		}
	})
	http.HandleFunc("/logout", authHandler.Logout)
	http.HandleFunc("/profile", authHandler.Profile)
	http.HandleFunc("/topic_create", topicHandler.CreateTopic)
	http.HandleFunc("/topic_new", topicHandler.ShowNewTopic)
	http.HandleFunc("/categories", categoryHandler.ShowCategories)
	http.HandleFunc("/topics", topicHandler.ShowTopics)
	http.HandleFunc("/topic", topicHandler.ShowTopic)
	http.HandleFunc("/post_create", postHandler.ShowCreatePostForm)
	http.HandleFunc("/post_create_submit", postHandler.CreatePost)

	// ─── Google OAuth Routes ──────────────────────────────────────────────────
	http.HandleFunc("/auth/google/login", authHandler.GoogleLogin)
	http.HandleFunc("/auth/google/callback", authHandler.GoogleCallback)

	// ─── Static Files ─────────────────────────────────────────────────────────
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	go middleware.CleanupVisitors()

	log.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", middleware.RateLimit(http.DefaultServeMux))
}
