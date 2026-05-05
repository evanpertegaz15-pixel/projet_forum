// création de commentaire, édition, suppression, permissions

package services

import (
	//"errors"
	"forum-dark-jurassic/internal/models"
)

type CommentService struct {
	Comments *models.CommentModel
	Posts    *models.PostModel
}

func NewCommentService(comments *models.CommentModel, posts *models.PostModel) *CommentService {
	return &CommentService{
		Comments: comments,
		Posts:    posts,
	}
}

/*
// Vérification admin
func isAdmin(user *models.User) bool {
	return user != nil && user.Role == "admin"
}

// Créer un commentaire
func (s *CommentService) CreateComment(user *models.User, postID int, content string) (int, error) {
	if user == nil {
		return 0, errors.New("utilisateur non connecté")
	}

	if content == "" {
		return 0, errors.New("contenu vide")
	}

	// Vérifier que le post existe
	post, err := s.Posts.FindByID(postID)
	if err != nil {
		return 0, err
	}
	if post == nil {
		return 0, errors.New("post introuvable")
	}

	return s.Comments.Create(postID, user.ID, content)
}

// Modifier un commentaire
func (s *CommentService) UpdateComment(user *models.User, commentID int, newContent string) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	comment, err := s.Comments.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("commentaire introuvable")
	}

	// Permission : auteur OU admin
	if comment.UserID != user.ID && !isAdmin(user) {
		return errors.New("permission refusée")
	}

	if newContent == "" {
		return errors.New("contenu vide")
	}

	return s.Comments.Update(commentID, newContent)
}

// Supprimer un commentaire
func (s *CommentService) DeleteComment(user *models.User, commentID int) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	comment, err := s.Comments.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("commentaire introuvable")
	}

	// Permission : auteur OU admin
	if comment.UserID != user.ID && !isAdmin(user) {
		return errors.New("permission refusée")
	}

	return s.Comments.Delete(commentID)
}

// Obtenir les commentaires d’un post
func (s *CommentService) GetCommentsByPost(postID int) ([]*models.CommentModel, error) {
	return s.Comments.GetByPostID(postID)
}
*/