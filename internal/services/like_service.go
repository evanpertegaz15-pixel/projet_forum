// like post, commentaire, unlike, règles (pas 2 likes du même type)

package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
)

type LikeService struct {
	Likes *models.LikeModel
	Posts *models.PostModel
	Comments *models.CommentModel
}

func NewLikeService(likes *models.LikeModel, posts *models.PostModel, comments *models.CommentModel) *LikeService {
	return &LikeService{
		Likes: likes,
		Posts: posts,
		Comments: comments,
	}
}

// Like un post
func (s *LikeService) LikePost(user *models.User, postID int) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	// vérifier post existe
	post, err := s.Posts.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post introuvable")
	}

	// vérifier si déjà liké
	exists, err := s.Likes.Exists(user.ID, postID, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("déjà liké")
	}

	return s.Likes.Create(user.ID, postID, 0)
}

// Unlike un post
func (s *LikeService) UnlikePost(user *models.User, postID int) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	return s.Likes.Delete(user.ID, postID, 0)
}

// Like un commentaire
func (s *LikeService) LikeComment(user *models.User, commentID int) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	// vérifier commentaire existe
	comment, err := s.Comments.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("commentaire introuvable")
	}

	// vérifier déjà liké
	exists, err := s.Likes.Exists(user.ID, 0, commentID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("déjà liké")
	}

	return s.Likes.Create(user.ID, 0, commentID)
}

// Unlike un commentaire
func (s *LikeService) UnlikeComment(user *models.User, commentID int) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	return s.Likes.Delete(user.ID, 0, commentID)
}

// Compter likes d’un post
func (s *LikeService) CountPostLikes(postID int) (int, error) {
	return s.Likes.Count(postID, 0)
}

// Compter likes d’un commentaire
func (s *LikeService) CountCommentLikes(commentID int) (int, error) {
	return s.Likes.Count(0, commentID)
}
