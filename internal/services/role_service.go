package services

import (
	"forum-dark-jurassic/internal/models"
)

type RoleService struct {
	Users *models.UserModel
}

func NewRoleService(users *models.UserModel) *RoleService {
	return &RoleService{
		Users: users,
	}
}

const (
	RoleUser    = "user"
	RoleAdmin   = "admin"
	RoleMod     = "moderator"
	RoleRanger  = "ranger"
	RoleBlocked = "blocked"
)

func (s *RoleService) IsAdmin(user *models.User) bool {
	return user != nil && user.HasRole(RoleAdmin)
}

func (s *RoleService) IsModerator(user *models.User) bool {
	return user != nil && (user.HasRole(RoleMod) || s.IsAdmin(user))
}

func (s *RoleService) IsRanger(user *models.User) bool {
	return user != nil && user.HasRole(RoleRanger)
}

func (s *RoleService) IsBlocked(user *models.User) bool {
	return user != nil && user.HasRole(RoleBlocked)
}

func (s *RoleService) CanCreateTopic(user *models.User) bool {
	return user != nil && !s.IsBlocked(user)
}

func (s *RoleService) CanCreatePost(user *models.User) bool {
	return user != nil && !s.IsBlocked(user)
}

func (s *RoleService) CanCreateReply(user *models.User) bool {
	return user != nil && !s.IsBlocked(user)
}

func (s *RoleService) CanModerateContent(user *models.User) bool {
	return s.IsModerator(user)
}

func (s *RoleService) CanManageReports(user *models.User) bool {
	return s.IsModerator(user)
}
