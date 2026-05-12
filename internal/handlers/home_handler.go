package handlers

import (
    "net/http"
    "forum-dark-jurassic/internal/services"
    "forum-dark-jurassic/internal/utils"
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
    utils.Render(w, "./internal/templates/home.html", map[string]any{
        "Users": users,
    })
}