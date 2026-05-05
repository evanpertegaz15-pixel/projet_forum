// assigner, retirer un rôle, vérifier les permissions

package services

import (
	//"errors"
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

// rôles possibles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
	RoleMod   = "moderator"
)

/*
// Vérifier admin
func (s *RoleService) IsAdmin(user *models.User) bool {
	return user != nil && user.Role == RoleAdmin
}

// Vérifier modérateur ou admin
func (s *RoleService) IsModerator(user *models.User) bool {
	return user != nil && (user.Role == RoleMod || user.Role == RoleAdmin)
}

// Assigner un rôle
func (s *RoleService) AssignRole(currentUser *models.User, targetUserID int, role string) error {
	if !s.IsAdmin(currentUser) {
		return errors.New("permission refusée")
	}

	// vérifier rôle valide
	if role != RoleUser && role != RoleAdmin && role != RoleMod {
		return errors.New("rôle invalide")
	}

	user, err := s.Users.FindByID(targetUserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("utilisateur introuvable")
	}

	return s.Users.UpdateRole(targetUserID, role)
}

// Retirer rôle (remettre en user)
func (s *RoleService) RemoveRole(currentUser *models.User, targetUserID int) error {
	if !s.IsAdmin(currentUser) {
		return errors.New("permission refusée")
	}

	user, err := s.Users.FindByID(targetUserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("utilisateur introuvable")
	}

	return s.Users.UpdateRole(targetUserID, RoleUser)
}

// Vérifier permission générique
func (s *RoleService) HasPermission(user *models.User, permission string) bool {
	if user == nil {
		return false
	}

	switch permission {

	case "manage_users":
		return user.Role == RoleAdmin

	case "moderate_content":
		return user.Role == RoleAdmin || user.Role == RoleMod

	case "create_post":
		return true

	case "delete_any_post":
		return user.Role == RoleAdmin

	default:
		return false
	}
}
*/