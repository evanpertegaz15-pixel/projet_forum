package handlers

import (
    "net/http"
    "strconv"
    "forum-dark-jurassic/internal/services"
)

type PostHandler struct {
    Posts *services.PostService
}

func NewPostHandler(p *services.PostService) *PostHandler {
    return &PostHandler{Posts: p}
}

func (handler *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée.", http.StatusMethodNotAllowed)
        return
    }
    topicID, _ := strconv.Atoi(r.FormValue("topic_id"))
    userID := 1 // temporaire -> session
    content := r.FormValue("content")
    _, err := handler.Posts.CreatePost(topicID, userID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/topic?id="+strconv.Itoa(topicID), http.StatusSeeOther)
}