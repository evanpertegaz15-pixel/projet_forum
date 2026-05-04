package models

import (
	"database/sql"
	"time"
)

type Like struct {
	ID        int
	UserID    int
	PostID    int
	CreatedAt time.Time
}

type LikeModel struct {
	DB *sql.DB
}

func NewLikeModel(db *sql.DB) *LikeModel {
	return &LikeModel{DB: db}
}

// Ajouter un like sur un post
func (m *LikeModel) AddLike(userID, postID int) error {
	// Empêche les doublons via logique SQL ou contrainte UNIQUE (user_id, post_id)
	query := `
		INSERT INTO likes (user_id, post_id, created_at)
		VALUES (?, ?, ?)
	`

	_, err := m.DB.Exec(query, userID, postID, time.Now())
	return err
}

// Supprimer un like (unlike)
func (m *LikeModel) RemoveLike(userID, postID int) error {
	query := `
		DELETE FROM likes
		WHERE user_id = ? AND post_id = ?
	`

	_, err := m.DB.Exec(query, userID, postID)
	return err
}

// Vérifier si un user a déjà liké un post
func (m *LikeModel) HasLiked(userID, postID int) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM likes
		WHERE user_id = ? AND post_id = ?
	`

	var count int
	err := m.DB.QueryRow(query, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Compter les likes d’un post
func (m *LikeModel) CountLikes(postID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM likes
		WHERE post_id = ?
	`

	var count int
	err := m.DB.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}