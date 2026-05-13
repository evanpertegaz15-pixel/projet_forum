package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/utils"
)

type ProfileHandler struct {
	Users *models.UserModel
}

func NewProfileHandler(userModel *models.UserModel) *ProfileHandler {
	return &ProfileHandler{
		Users: userModel,
	}
}

func (h *ProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		utils.ErrorBadRequest(w, "id utilisateur invalide")
		return
	}
	user, err := h.Users.FindByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		utils.ErrorNotFound(w, "Utilisateur introuvable")
		return
	}
	response := map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"profile_picture":	user.ProfilePicture,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		utils.ErrorUnauthorized(w, "Non authentifié")
		return
	}
	user, err := h.Users.FindByID(userID.(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		utils.ErrorNotFound(w, "Utilisateur introuvable")
		return
	}
	response := map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"profile_picture":	user.ProfilePicture,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}