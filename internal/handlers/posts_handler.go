package handlers

import (
    "html/template"
    "net/http"
    "strconv"
    "forum-dark-jurassic/internal/services"
)

type PostHandler struct {
    Posts *services.PostService
    Auth *services.AuthService
}

func NewPostHandler(posts *services.PostService, auth *services.AuthService) *PostHandler {
    return &PostHandler{Posts: posts, Auth: auth,}
}

func (handler *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée.", http.StatusMethodNotAllowed)
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
    topicID, err := strconv.Atoi(r.FormValue("topic_id"))
    if err != nil || topicID <= 0 {
        http.Error(w, "Topic invalide.", http.StatusBadRequest)
        return
    }
    content := r.FormValue("content")
    if content == "" {
        http.Error(w, "Le contenu ne peut pas être vide.", http.StatusBadRequest)
        return
    }
    _, err = handler.Posts.CreatePost(topicID, user.ID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/topic?id="+strconv.Itoa(topicID), http.StatusSeeOther)
}

func (handler *PostHandler) ShowCreatePostForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Méthode non autorisée.", http.StatusMethodNotAllowed)
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
    tmpl := template.Must(template.ParseFiles("./internal/templates/post_create.html"))
    tmpl.Execute(w, user)
}