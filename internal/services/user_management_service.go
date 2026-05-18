package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
)

type UserManagementService struct {
	Users *models.UserModel
	UserRoles *models.UserRoleModel
	Roles *models.RoleModel
}

func NewUserManagementService(users *models.UserModel, userRoles *models.UserRoleModel, roles *models.RoleModel) *UserManagementService {
	return &UserManagementService{
		Users: users,
		UserRoles: userRoles,
		Roles: roles,
	}
}

func (service *UserManagementService) GetAllUsers() ([]models.User, error) {
	return service.Users.GetAllUsers()
}

func (service *UserManagementService) AssignRole(userID, roleID int) error {
	return service.UserRoles.AssignRole(userID, roleID)
}

func (service *UserManagementService) RemoveRole(userID, roleID int) error {
	return service.UserRoles.RemoveRole(userID, roleID)
}

func (service *UserManagementService) GetAllRoles() ([]models.Role, error) {
	return service.Roles.GetAllRoles()
}

func (service *UserManagementService) BlockUser(userID int) error {
	blockedRole, err := service.Roles.GetRoleByName("blocked")
	if err != nil {
		return err
	}
	if blockedRole == nil {
		return errors.New("rôle 'blocked' introuvable")
	}
	return service.UserRoles.AssignRole(userID, blockedRole.ID)
}

func (service *UserManagementService) UnblockUser(userID int) error {
	blockedRole, err := service.Roles.GetRoleByName("blocked")
	if err != nil {
		return err
	}
	if blockedRole == nil {
		return errors.New("rôle 'blocked' introuvable")
	}
	return service.UserRoles.RemoveRole(userID, blockedRole.ID)
}

func (service *UserManagementService) PromoteToModerator(userID int) error {
	modRole, err := service.Roles.GetRoleByName("moderator")
	if err != nil {
		return err
	}
	if modRole == nil {
		return errors.New("rôle 'moderator' introuvable")
	}
	return service.UserRoles.AssignRole(userID, modRole.ID)
}

func (service *UserManagementService) PromoteToRanger(userID int) error {
	rangerRole, err := service.Roles.GetRoleByName("ranger")
	if err != nil {
		return err
	}
	if rangerRole == nil {
		return errors.New("rôle 'ranger' introuvable")
	}
	return service.UserRoles.AssignRole(userID, rangerRole.ID)
}

func (service *UserManagementService) DemoteFromModerator(userID int) error {
	modRole, err := service.Roles.GetRoleByName("moderator")
	if err != nil {
		return err
	}
	if modRole == nil {
		return errors.New("rôle 'moderator' introuvable")
	}
	return service.UserRoles.RemoveRole(userID, modRole.ID)
}

func (service *UserManagementService) DemoteFromRanger(userID int) error {
	rangerRole, err := service.Roles.GetRoleByName("ranger")
	if err != nil {
		return err
	}
	if rangerRole == nil {
		return errors.New("rôle 'ranger' introuvable")
	}
	return service.UserRoles.RemoveRole(userID, rangerRole.ID)
}
