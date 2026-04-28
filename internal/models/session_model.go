package models

import (
	"database/sql"
	"errors"
	"time"
	"forum-dark-jurassic/internal/utils"
)

type Session struct {
	ID	string
	UserID	int
	CreatedAt	time.Time
	ExpiresAt	time.Time
}

type SessionModel struct {
	DB *sql.DB
}

func NewSessionModel(db *sql.DB) *SessionModel {
	return &SessionModel{DB: db}
}

func (model *SessionModel) CreateSession(userID int) (string, error) {
	id := utils.NewUUID()
	expires := time.Now().Add(7 * 24 * time.Hour) // = 7 jours
	_, err := model.DB.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)
	`, id, userID, expires)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (model *SessionModel) GetSession(id string) (*Session, error) {
	row := model.DB.QueryRow(`
		SELECT id, user_id, created_at, expires_at
		FROM sessions
		WHERE id = ?
	`, id)
	var session Session
	err := row.Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (model *SessionModel) DeleteSession(id string) error {
	_, err := model.DB.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

func (model *SessionModel) DeleteSessionByUser(userID int) error {
	_, err := model.DB.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	return err
}