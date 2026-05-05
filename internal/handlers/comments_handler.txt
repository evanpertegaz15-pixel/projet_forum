package handlers
// comments -> ajouter, modifier, supprimer un commentaire
import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum-dark-jurassic/internal/services"
)

type CommentHandler struct {
	Comments *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		Comments: commentService,
	}
}

// Ajouter un commentaire
func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(int)

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "post_id invalide", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "commentaire vide", http.StatusBadRequest)
		return
	}

	err = h.Comments.CreateComment(content, userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Commentaire ajouté",
	})
}

// Modifier un commentaire
func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(int)

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "comment_id invalide", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "contenu vide", http.StatusBadRequest)
		return
	}

	err = h.Comments.UpdateComment(commentID, userID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Commentaire modifié",
	})
}

// Supprimer un commentaire
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(int)

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "comment_id invalide", http.StatusBadRequest)
		return
	}

	err = h.Comments.DeleteComment(commentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Commentaire supprimé",
	})
}