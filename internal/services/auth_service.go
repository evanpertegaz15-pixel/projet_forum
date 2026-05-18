package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/utils"
	"time"
)

type AuthService struct {
	Users    *models.UserModel
	Sessions *models.SessionModel
}

func NewAuthService(users *models.UserModel, sessions *models.SessionModel) *AuthService {
	return &AuthService{
		Users:    users,
		Sessions: sessions,
	}
}

func (service *AuthService) CheckPassword(user *models.User, password string) bool {
    return utils.CheckPasswordHash(user.Password, password)
}

func (service *AuthService) CreateSession(userID int) (string, error) {
	id := utils.NewUUID()
	expires := time.Now().Add(7 * 24 * time.Hour) // = 7 jours
	err := service.Sessions.InsertSession(id, userID, expires)
	if err != nil {
		return "", err
	}
	return id, nil
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
	if err != nil {
		return 0, err
	}
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
	isPasswordValid := auth.CheckPassword(user, password)
	if !isPasswordValid {
		return "", errors.New("Mot de passe incorrect.")
	}
	return auth.CreateSession(user.ID)
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

// ─── Google OAuth ─────────────────────────────────────────────────────────────

// FindOrCreateGoogleUser looks up a user by email.
// If found, returns their ID. If not, creates a new one with an empty password.
func (auth *AuthService) FindOrCreateGoogleUser(email, name string) (int64, error) {
	user, err := auth.Users.FindByEmail(email)
	if err != nil {
		return 0, err
	}
	if user != nil {
		return int64(user.ID), nil
	}
	id, err := auth.Users.CreateUser(email, name, "")
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

// CreateSession exposes session creation so GoogleCallback can call it directly.
func (auth *AuthService) CreateSessionFromGoogle(userID int64) (string, error) {
	return auth.CreateSession(int(userID))
}

func (auth *AuthService) UserHasRole(user *models.User, role string) bool {
    for _, r := range user.Roles {
        if r.Name == role {
            return true
        }
    }
    return false
}

func (auth *AuthService) UpdateEmail(userID int, newEmail string) error {
	return auth.Users.UpdateEmail(userID, newEmail)
}

func (auth *AuthService) UpdateUsername(userID int, newUsername string) error {
	return auth.Users.UpdateUsername(userID, newUsername)
}

func (auth *AuthService) UpdatePassword(userID int, newPassword string) error {
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return auth.Users.UpdatePassword(userID, hashed)
}
