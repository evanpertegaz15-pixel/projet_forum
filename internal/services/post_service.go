package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
)

type PostService struct {
	Posts       *models.PostModel
}

func NewPostService(posts *models.PostModel) *PostService {
	return &PostService{Posts: posts}
}

func (service *PostService) GetPostsByTopic(topicID int) ([]models.Post, error) {
    return service.Posts.GetPostsByTopic(topicID)
}

func (service *PostService) CreatePost(topicID, userID int, content string) (int, error) {
    if content == "" {
        return 0, errors.New("Le contenu ne peut pas être vide.")
    }
    return service.Posts.CreatePost(topicID, userID, content)
}

func (service *PostService) CreateReply(topicID, userID, parentID int, content string) (int, error) {
    if content == "" {
        return 0, errors.New("Le contenu ne peut pas être vide.")
    }
    return service.Posts.CreateReply(topicID, userID, parentID, content)
}

func (service *PostService) GetReplies(postID int) ([]models.Post, error) {
    return service.Posts.GetReplies(postID)
}

func (service *PostService) GetPostsWithRepliesByTopic(topicID int) ([]models.PostWithReplies, error) {
    return service.Posts.GetPostsWithRepliesByTopic(topicID)
}

func (s *PostService) DeletePost(user *models.User, postID int) error {
	post, err := s.Posts.GetPostByID(postID)
	if err != nil {
		return err
	}
	if post.ID == 0 {
		return errors.New("post introuvable")
	}
	if post.UserID != user.ID && !user.HasRole("admin") && !user.HasRole("moderator") && !user.HasRole("ranger") {
		return errors.New("permission refusée")
	}
	return s.Posts.Delete(postID)
}