package handlers

import (
    "html/template"
    "net/http"
    "strconv"
    "forum-dark-jurassic/internal/services"
)

type TopicHandler struct {
    Topics *services.TopicService
    Posts  *services.PostService
}

func NewTopicHandler(topics *services.TopicService, posts *services.PostService) *TopicHandler {
    return &TopicHandler{
        Topics: topics,
        Posts:  posts,
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
    posts, err := handler.Posts.GetPostsByTopic(topicID)
    if err != nil {
        http.Error(w, "Erreur interne.", http.StatusInternalServerError)
        return
    }
    tmpl := template.Must(template.ParseFiles("./internal/templates/topic.html"))
    tmpl.Execute(w, posts)
}

func (handler *TopicHandler) ShowNewTopic(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("./internal/templates/new_topic.html"))
    tmpl.Execute(w, nil)
}

func (handler *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }
    categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
    userID := 1 // temporaire -> session
    title := r.FormValue("title")
    _, err := handler.Topics.CreateTopic(categoryID, userID, title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/category?id="+strconv.Itoa(categoryID), http.StatusSeeOther)
}