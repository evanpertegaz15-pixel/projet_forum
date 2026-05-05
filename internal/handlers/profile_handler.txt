package handlers
// users -> afficher un profil utilisateur
import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum-dark-jurassic/internal/models"
)

type ProfileHandler struct {
	Users *models.UserModel
}

func NewProfileHandler(userModel *models.UserModel) *ProfileHandler {
	return &ProfileHandler{
		Users: userModel,
	}
}

// Afficher le profil d’un utilisateur par ID
func (h *ProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("id")

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		http.Error(w, "id utilisateur invalide", http.StatusBadRequest)
		return
	}

	user, err := h.Users.FindByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	// On évite de renvoyer le hash du mot de passe
	response := map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Afficher le profil de l’utilisateur connecté
func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	user, err := h.Users.FindByID(userID.(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	response := map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}