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

	reportModel := models.NewReportModel(db)
	userModel := models.NewUserModel(db)
	sessionModel := models.NewSessionModel(db)
	postModel := models.NewPostModel(db)
	categoryModel := models.NewCategoryModel(db)
	topicModel := models.NewTopicModel(db)
	likeModel := models.NewLikeModel(db)

	userService := services.NewUserService(userModel)
	authService := services.NewAuthService(userModel, sessionModel)
	postService := services.NewPostService(postModel)
	categoryService := services.NewCategoryService(categoryModel)
	topicService := services.NewTopicService(topicModel)
	likeService := services.NewLikeService(likeModel, postModel, topicModel)
	reportService := services.NewReportService(reportModel, postModel, topicModel, userModel)

	homeHandler := handlers.NewHomeHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	topicHandler := handlers.NewTopicHandler(topicService, postService, categoryService, authService)
	postHandler := handlers.NewPostHandler(postService, authService)
	likesHandler := handlers.NewLikesHandler(likeService, authService)
	reportHandler := handlers.NewReportHandler(reportService, authService)

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
	http.HandleFunc("/delete/account", authHandler.DeleteAccount)
	http.HandleFunc("/topic/create", topicHandler.CreateTopic)
	http.HandleFunc("/topic/new", topicHandler.ShowNewTopic)
	http.HandleFunc("/categories", categoryHandler.ShowCategories)
	http.HandleFunc("/topics", topicHandler.ShowTopics)
	http.HandleFunc("/topic", topicHandler.ShowTopic)
	http.HandleFunc("/post/create", postHandler.ShowCreatePostForm)
	http.HandleFunc("/post/create/submit", postHandler.CreatePost)
	http.HandleFunc("/post", postHandler.ShowPost)
	http.HandleFunc("/reply", postHandler.CreateReply)
	http.HandleFunc("/post/delete", postHandler.DeletePost)
	http.HandleFunc("/topic/delete", topicHandler.DeleteTopic)
	http.HandleFunc("/report", reportHandler.CreateReport)
	http.HandleFunc("/reports", reportHandler.ShowReports)
	http.HandleFunc("/reports/open", reportHandler.GetOpenReports)
	http.HandleFunc("/report/resolve", reportHandler.ResolveReport)
	http.HandleFunc("/report/delete", reportHandler.DeleteReport)
	http.HandleFunc("/report/delete-content", reportHandler.DeleteReportedContent)
	
	// ─── Google OAuth Routes ──────────────────────────────────────────────────
	http.HandleFunc("/auth/google/login", authHandler.GoogleLogin)
	http.HandleFunc("/auth/google/callback", authHandler.GoogleCallback)

	// ─── Static Files ─────────────────────────────────────────────────────────
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	go middleware.CleanupVisitors()

	http.HandleFunc("/like/post", likesHandler.TogglePostLike)
	http.HandleFunc("/like/topic", likesHandler.ToggleTopicLike)
	http.HandleFunc("/like/comment", likesHandler.ToggleCommentLike)
	http.HandleFunc("/likes/post", likesHandler.GetPostLikesCount)
	http.HandleFunc("/likes/topic", likesHandler.GetTopicLikesCount)
	http.HandleFunc("/likes/comment", likesHandler.GetCommentLikesCount)

	log.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", middleware.RateLimit(http.DefaultServeMux))
}
