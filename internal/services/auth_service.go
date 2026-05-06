package services

import (
	"time"
	"errors"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/utils"
)

type AuthService struct {
	Users	*models.UserModel
	Sessions	*models.SessionModel
}

func NewAuthService(users *models.UserModel, sessions *models.SessionModel) *AuthService {
	return &AuthService{
		Users:	users,
		Sessions:	sessions,
	}
}

func (auth *AuthService) Register(email, username, password string) (int, error) {
	existing, err := auth.Users.FindByEmail(email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("Cet email est déjà utilisé.")
	}
	hashed, err := utils.HashPassword(password)
	return auth.Users.CreateUser(email, username, hashed)
}

func (auth *AuthService) Login(identifier, password string) (string, error) {
	user, err := auth.Users.FindByEmail(identifier)
	if err != nil {
		return "", err
	}
	if user == nil {
		user, err = auth.Users.FindByUsername(identifier)
		if err != nil {
			return "", err
		}
	}
	if user == nil {
		return "", errors.New("Utilisateur non trouvé.")
	}
	
	// DEBUG: Vérify password hash
	isPasswordValid := user.CheckPassword(password)
	if !isPasswordValid {
		return "", errors.New("Mot de passe incorrect.")
	}
	return auth.Sessions.CreateSession(user.ID)
}

func (auth *AuthService) Logout(sessionID string) error {
	return auth.Sessions.DeleteSession(sessionID)
}

func (auth *AuthService) GetUserFromSession(sessionID string) (*models.User, error) {
	session, err := auth.Sessions.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	if time.Now().After(session.ExpiresAt) {
		_ = auth.Sessions.DeleteSession(sessionID)
		return nil, nil
	}
	return auth.Users.FindByID(session.UserID)
}