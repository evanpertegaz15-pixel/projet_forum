package handlers
// likes -> like / unlike post ou commentaire
import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum-dark-jurassic/internal/services"
)

type LikesHandler struct {
	Likes *services.LikeService
}

func NewLikesHandler(likeService *services.LikeService) *LikesHandler {
	return &LikesHandler{
		Likes: likeService,
	}
}

// Like / Unlike un post (toggle)
func (h *LikesHandler) TogglePostLike(w http.ResponseWriter, r *http.Request) {
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

	liked, err := h.Likes.TogglePostLike(userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"liked": liked,
	})
}

// Like / Unlike un commentaire (toggle)
func (h *LikesHandler) ToggleCommentLike(w http.ResponseWriter, r *http.Request) {
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

	liked, err := h.Likes.ToggleCommentLike(userID, commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"liked": liked,
	})
}

// Récupérer le nombre de likes d’un post
func (h *LikesHandler) GetPostLikesCount(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "post_id invalide", http.StatusBadRequest)
		return
	}

	count, err := h.Likes.CountPostLikes(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}

// Récupérer le nombre de likes d’un commentaire
func (h *LikesHandler) GetCommentLikesCount(w http.ResponseWriter, r *http.Request) {
	commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "comment_id invalide", http.StatusBadRequest)
		return
	}

	count, err := h.Likes.CountCommentLikes(commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}