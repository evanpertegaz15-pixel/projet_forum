package models
// associer / dissocier un tag d'un post
import (
	"database/sql"
)

type PostTag struct {
	PostID int
	TagID  int
}

type PostTagModel struct {
	DB *sql.DB
}

func NewPostTagModel(db *sql.DB) *PostTagModel {
	return &PostTagModel{DB: db}
}

// Associer un tag à un post
func (m *PostTagModel) AddTagToPost(postID, tagID int) error {
	query := `
		INSERT INTO post_tags (post_id, tag_id)
		VALUES (?, ?)
	`

	_, err := m.DB.Exec(query, postID, tagID)
	return err
}

// Dissocier un tag d'un post
func (m *PostTagModel) RemoveTagFromPost(postID, tagID int) error {
	query := `
		DELETE FROM post_tags
		WHERE post_id = ? AND tag_id = ?
	`

	_, err := m.DB.Exec(query, postID, tagID)
	return err
}

// Supprimer tous les tags d’un post
func (m *PostTagModel) ClearTagsFromPost(postID int) error {
	query := `
		DELETE FROM post_tags
		WHERE post_id = ?
	`

	_, err := m.DB.Exec(query, postID)
	return err
}

// Récupérer tous les tags d’un post
func (m *PostTagModel) GetTagsByPostID(postID int) ([]int, error) {
	query := `
		SELECT tag_id
		FROM post_tags
		WHERE post_id = ?
	`

	rows, err := m.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []int

	for rows.Next() {
		var tagID int
		err := rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tagID)
	}

	return tags, nil
}

// Récupérer tous les posts d’un tag
func (m *PostTagModel) GetPostsByTagID(tagID int) ([]int, error) {
	query := `
		SELECT post_id
		FROM post_tags
		WHERE tag_id = ?
	`

	rows, err := m.DB.Query(query, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []int

	for rows.Next() {
		var postID int
		err := rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, postID)
	}

	return posts, nil
}