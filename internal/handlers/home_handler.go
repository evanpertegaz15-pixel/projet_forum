package handlers

import (
    "html/template"
    "net/http"
    "forum-dark-jurassic/internal/services"
)

type HomeHandler struct {
    Users *services.UserService
}

func NewHomeHandler(users *services.UserService) *HomeHandler {
    return &HomeHandler{Users: users}
}

func (handler *HomeHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
    users, err := handler.Users.GetAllUsers()
    if err != nil {
        http.Error(w, "Erreur interne.", http.StatusInternalServerError)
        return
    }
    tmpl := template.Must(template.ParseFiles("./internal/templates/home.html"))
    tmpl.Execute(w, map[string]any{
        "Users": users,
    })
}