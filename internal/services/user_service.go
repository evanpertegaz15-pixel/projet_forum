package services

import "forum-dark-jurassic/internal/models"

type UserService struct {
    Users *models.UserModel
}

func NewUserService(users *models.UserModel) *UserService {
    return &UserService{Users: users}
}

func (service *UserService) GetAllUsers() ([]models.User, error) {
    return service.Users.GetAllUsers()
}