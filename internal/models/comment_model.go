package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	Content   string
	UserID    int
	PostID    int
	CreatedAt time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func NewCommentModel(db *sql.DB) *CommentModel {
	return &CommentModel{DB: db}
}

// Créer un commentaire
func (m *CommentModel) CreateComment(content string, userID, postID int) (int, error) {
	query := `
		INSERT INTO comments (content, user_id, post_id, created_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := m.DB.Exec(query, content, userID, postID, time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Récupérer les commentaires d'un post
func (m *CommentModel) GetCommentsByPostID(postID int) ([]Comment, error) {
	query := `
		SELECT id, content, user_id, post_id, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at ASC
	`

	rows, err := m.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.Content, &c.UserID, &c.PostID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// Supprimer un commentaire (optionnel mais utile)
func (m *CommentModel) DeleteComment(commentID int, userID int) error {
	query := `
		DELETE FROM comments
		WHERE id = ? AND user_id = ?
	`

	_, err := m.DB.Exec(query, commentID, userID)
	return err
}