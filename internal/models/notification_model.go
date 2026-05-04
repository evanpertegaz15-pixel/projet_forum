package models
// lister, marquer comme lu, supprimer les anciennespackage models

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID        int
	UserID    int
	Message   string
	IsRead    bool
	CreatedAt time.Time
}

type NotificationModel struct {
	DB *sql.DB
}

func NewNotificationModel(db *sql.DB) *NotificationModel {
	return &NotificationModel{DB: db}
}

// Créer une notification
func (m *NotificationModel) CreateNotification(userID int, message string) error {
	query := `
		INSERT INTO notifications (user_id, message, is_read, created_at)
		VALUES (?, ?, 0, ?)
	`

	_, err := m.DB.Exec(query, userID, message, time.Now())
	return err
}

// Lister les notifications d’un user (plus récentes en premier)
func (m *NotificationModel) GetNotifications(userID int) ([]Notification, error) {
	query := `
		SELECT id, user_id, message, is_read, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var n Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// Marquer une notification comme lue
func (m *NotificationModel) MarkAsRead(notificationID int, userID int) error {
	query := `
		UPDATE notifications
		SET is_read = 1
		WHERE id = ? AND user_id = ?
	`

	_, err := m.DB.Exec(query, notificationID, userID)
	return err
}

// Marquer toutes les notifications comme lues
func (m *NotificationModel) MarkAllAsRead(userID int) error {
	query := `
		UPDATE notifications
		SET is_read = 1
		WHERE user_id = ?
	`

	_, err := m.DB.Exec(query, userID)
	return err
}

// Supprimer les anciennes notifications (ex: +30 jours)
func (m *NotificationModel) DeleteOldNotifications(days int) error {
	query := `
		DELETE FROM notifications
		WHERE created_at < datetime('now', ?)
	`

	// SQLite syntax : "-30 days"
	interval := "-" + time.Duration(days*24).String() + "h"

	_, err := m.DB.Exec(query, interval)
	return err
}

// Supprimer une notification précise
func (m *NotificationModel) DeleteNotification(notificationID int, userID int) error {
	query := `
		DELETE FROM notifications
		WHERE id = ? AND user_id = ?
	`

	_, err := m.DB.Exec(query, notificationID, userID)
	return err
}