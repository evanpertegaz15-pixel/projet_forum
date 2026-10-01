package handlers

import (
	"net/http"
	"forum-dark-jurassic/internal/utils"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (handler *NotFoundHandler) ShowNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	utils.Render(w, "./internal/templates/404.html", nil)
}