package handlers

import (
    "html/template"
    "log"
    "net/http"
    "strconv"
    "forum-dark-jurassic/internal/models"
    "forum-dark-jurassic/internal/services"
)

type TopicHandler struct {
    Topics *services.TopicService
    Posts  *services.PostService
    Categories *services.CategoryService
    Auth *services.AuthService
}

func NewTopicHandler(topics *services.TopicService, posts *services.PostService, categories *services.CategoryService, auth *services.AuthService) *TopicHandler {
    return &TopicHandler{
        Topics: topics,
        Posts:  posts,
        Categories: categories,
        Auth:   auth,
    }
}

func (handler *TopicHandler) ShowTopics(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    categoryID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Catégorie incorrecte.", http.StatusBadRequest)
        return
    }
    topics, err := handler.Topics.GetTopicsByCategory(categoryID)
    if err != nil {
        http.Error(w, "Erreur interne.", http.StatusInternalServerError)
        return
    }
    tmpl := template.Must(template.ParseFiles("./internal/templates/topics.html"))
    tmpl.Execute(w, topics)
}

func (handler *TopicHandler) ShowTopic(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    topicID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Topic incorrect.", http.StatusBadRequest)
        return
    }
    topic, err := handler.Topics.GetTopicByID(topicID)
    if err != nil {
        http.Error(w, "Topic introuvable.", http.StatusNotFound)
        return
    }
    posts, err := handler.Posts.GetPostsByTopic(topicID)
    if err != nil {
        http.Error(w, "Erreur interne.", http.StatusInternalServerError)
        return
    }
    data := struct {
        Topic models.Topic
        Posts []models.Post
    }{
        Topic: topic,
        Posts: posts,
    }
    tmpl := template.Must(template.ParseFiles("./internal/templates/topic.html"))
    tmpl.Execute(w, data)
}

func (handler *TopicHandler) ShowNewTopic(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    user, err := handler.Auth.GetUserFromSession(cookie.Value)
    if err != nil || user == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    categories, err := handler.Categories.GetAllCategories()
    if err != nil {
        log.Printf("ShowNewTopic: GetAllCategories error: %v", err)
        http.Error(w, "Impossible de charger les catégories.", http.StatusInternalServerError)
        return
    }
    tmpl := template.Must(template.ParseFiles("./internal/templates/new_topic.html"))
    if err := tmpl.Execute(w, map[string]any{
        "User": user,
        "Categories":   categories,
    }); err != nil {
        log.Printf("ShowNewTopic: template execute error: %v", err)
        http.Error(w, "Impossible de charger les catégories.", http.StatusInternalServerError)
        return
    }
}

func (handler *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }
    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    user, err := handler.Auth.GetUserFromSession(cookie.Value)
    if err != nil || user == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
    title := r.FormValue("title")
    content := r.FormValue("content")
    if categoryID <= 0 || title == "" || content == "" {
        http.Error(w, "Champs incorrects.", http.StatusBadRequest)
        return
    }
    topicID, err := handler.Topics.CreateTopic(categoryID, user.ID, title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    _, err = handler.Posts.CreatePost(topicID, user.ID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/topic?id="+strconv.Itoa(topicID), http.StatusSeeOther)
}