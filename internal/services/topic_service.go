package services

import "forum-dark-jurassic/internal/models"

type TopicService struct {
    Topics *models.TopicModel
}

func NewTopicService(topics *models.TopicModel) *TopicService {
    return &TopicService{Topics: topics}
}

func (service *TopicService) GetTopicsByCategory(categoryID int) ([]models.Topic, error) {
    return service.Topics.GetTopicByCategory(categoryID)
}

func (service *TopicService) CreateTopic(categoryID, userID int, title string) (int, error) {
    if title == "" {
        return 0, errors.New("Le titre ne peut pas être vide.")
    }
    return service.Topics.Create(categoryID, userID, title)
}