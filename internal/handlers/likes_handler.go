package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"forum-dark-jurassic/internal/services"
)

type LikesHandler struct {
	Likes *services.LikeService
	Auth  *services.AuthService
}

func NewLikesHandler(likeService *services.LikeService, authService *services.AuthService) *LikesHandler {
	return &LikesHandler{
		Likes: likeService,
		Auth:  authService,
	}
}

func (handler *LikesHandler) TogglePostLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	user, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "post_id invalide", http.StatusBadRequest)
		return
	}
	liked, err := handler.Likes.TogglePostLike(user, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if returnURL := r.FormValue("return_url"); returnURL != "" {
		http.Redirect(w, r, returnURL, http.StatusSeeOther)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"liked": liked,
	})
}

func (handler *LikesHandler) ToggleTopicLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	user, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	topicID, err := strconv.Atoi(r.FormValue("topic_id"))
	if err != nil {
		http.Error(w, "topic_id invalide", http.StatusBadRequest)
		return
	}
	liked, err := handler.Likes.ToggleTopicLike(user, topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if returnURL := r.FormValue("return_url"); returnURL != "" {
		http.Redirect(w, r, returnURL, http.StatusSeeOther)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"liked": liked,
	})
}

func (handler *LikesHandler) ToggleCommentLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	user, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "comment_id invalide", http.StatusBadRequest)
		return
	}
	liked, err := handler.Likes.ToggleCommentLike(user, commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if returnURL := r.FormValue("return_url"); returnURL != "" {
		http.Redirect(w, r, returnURL, http.StatusSeeOther)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"liked": liked,
	})
}

func (handler *LikesHandler) GetPostLikesCount(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "post_id invalide", http.StatusBadRequest)
		return
	}
	count, err := handler.Likes.CountPostLikes(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}

func (handler *LikesHandler) GetTopicLikesCount(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(r.URL.Query().Get("topic_id"))
	if err != nil {
		http.Error(w, "topic_id invalide", http.StatusBadRequest)
		return
	}
	count, err := handler.Likes.CountTopicLikes(topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}

func (handler *LikesHandler) GetCommentLikesCount(w http.ResponseWriter, r *http.Request) {
	commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "comment_id invalide", http.StatusBadRequest)
		return
	}
	count, err := handler.Likes.CountCommentLikes(commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}