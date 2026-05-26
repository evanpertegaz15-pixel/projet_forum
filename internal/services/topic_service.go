package services

import (
    "errors"
    "forum-dark-jurassic/internal/models")

type TopicService struct {
    Topics *models.TopicModel
}

func NewTopicService(topics *models.TopicModel) *TopicService {
    return &TopicService{Topics: topics}
}

func (service *TopicService) GetTopicByID(id int) (models.Topic, error) {
    return service.Topics.GetTopicByID(id)
}

func (service *TopicService) GetLatestTopics(limit int) ([]models.Topic, error) {
    return service.Topics.GetLatestTopics(limit)
}

func (service *TopicService) GetTopicsByCategory(categoryID int) ([]models.Topic, error) {
    return service.Topics.GetTopicsByCategory(categoryID)
}

func (service *TopicService) CreateTopic(categoryID, userID int, title string) (int, error) {
    if title == "" {
        return 0, errors.New("Le titre ne peut pas être vide.")
    }
    return service.Topics.CreateTopic(categoryID, userID, title)
}

func (service *TopicService) DeleteTopic(user *models.User, topicID int) error {
    topic, err := service.Topics.GetTopicByID(topicID)
    if err != nil {
        return err
    }
    if topic.ID == 0 {
        return errors.New("topic introuvable")
    }
    if topic.UserID != user.ID && !user.HasRole("admin") && !user.HasRole("moderator") && !user.HasRole("ranger") {
        return errors.New("permission refusée")
    }
    return service.Topics.DeleteTopic(topicID)
}

func (service *TopicService) GetLikedTopicsByUser(userID, categoryID int) ([]models.Topic, error) {
    return service.Topics.GetLikedTopicsByUser(userID, categoryID)
}