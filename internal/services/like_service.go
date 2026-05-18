package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
)

type LikeService struct {
	Likes  *models.LikeModel
	Posts  *models.PostModel
	Topics *models.TopicModel
}

func NewLikeService(likes *models.LikeModel, posts *models.PostModel, topics *models.TopicModel) *LikeService {
	return &LikeService{
		Likes:  likes,
		Posts:  posts,
		Topics: topics,
	}
}

func (service *LikeService) TogglePostLike(user *models.User, postID int) (bool, error) {
	if user == nil {
		return false, errors.New("utilisateur non connecté")
	}
	if postID <= 0 {
		return false, errors.New("post invalide")
	}
	_, err := service.Posts.GetPostByID(postID)
	if err != nil {
		return false, err
	}
	exists, err := service.Likes.HasLiked(user.ID, postID, 0, 0)
	if err != nil {
		return false, err
	}
	if exists {
		err = service.Likes.DeleteLike(user.ID, postID, 0, 0)
		return false, err
	}
	err = service.Likes.CreateLike(user.ID, postID, 0, 0)
	return err == nil, err
}

func (service *LikeService) ToggleTopicLike(user *models.User, topicID int) (bool, error) {
	if user == nil {
		return false, errors.New("utilisateur non connecté")
	}
	if topicID <= 0 {
		return false, errors.New("topic invalide")
	}
	_, err := service.Topics.GetTopicByID(topicID)
	if err != nil {
		return false, err
	}
	exists, err := service.Likes.HasLiked(user.ID, 0, 0, topicID)
	if err != nil {
		return false, err
	}
	if exists {
		err = service.Likes.DeleteLike(user.ID, 0, 0, topicID)
		return false, err
	}
	err = service.Likes.CreateLike(user.ID, 0, 0, topicID)
	return err == nil, err
}

func (service *LikeService) ToggleCommentLike(user *models.User, commentID int) (bool, error) {
	if user == nil {
		return false, errors.New("utilisateur non connecté")
	}
	if commentID <= 0 {
		return false, errors.New("commentaire invalide")
	}
	exists, err := service.Likes.HasLiked(user.ID, 0, commentID, 0)
	if err != nil {
		return false, err
	}
	if exists {
		err = service.Likes.DeleteLike(user.ID, 0, commentID, 0)
		return false, err
	}
	err = service.Likes.CreateLike(user.ID, 0, commentID, 0)
	return err == nil, err
}

func (s *LikeService) CountPostLikes(postID int) (int, error) {
	return s.Likes.CountLikes(postID, 0, 0)
}

func (s *LikeService) CountTopicLikes(topicID int) (int, error) {
	return s.Likes.CountLikes(0, 0, topicID)
}

func (s *LikeService) CountCommentLikes(commentID int) (int, error) {
	return s.Likes.CountLikes(0, commentID, 0)
}