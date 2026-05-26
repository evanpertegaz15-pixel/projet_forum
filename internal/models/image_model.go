package models
import (
	"database/sql"
	"time"
)

type Image struct {
	ID        int
	URL       string
	UserID    int
	PostID    *int // peut être NULL si image de profil par exemple
	CreatedAt time.Time
}

type ImageModel struct {
	DB *sql.DB
}

func NewImageModel(db *sql.DB) *ImageModel {
	return &ImageModel{DB: db}
}

func (m *ImageModel) CreateImage(url string, userID int, postID *int) (int, error) {
	query := `
		INSERT INTO images (path, user_id, post_id, created_at)
		VALUES (?, ?, ?, ?)
	`
	result, err := m.DB.Exec(query, url, userID, postID, time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ImageModel) GetImagesByPostID(postID int) ([]Image, error) {
	query := `
		SELECT id, path AS url, user_id, post_id, created_at
		FROM images
		WHERE post_id = ?
	`
	rows, err := m.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []Image
	for rows.Next() {
		var img Image
		err := rows.Scan(&img.ID, &img.URL, &img.UserID, &img.PostID, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func (m *ImageModel) GetImagesByUserID(userID int) ([]Image, error) {
	query := `
		SELECT id, path AS url, user_id, post_id, created_at
		FROM images
		WHERE user_id = ?
	`
	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []Image
	for rows.Next() {
		var img Image
		err := rows.Scan(&img.ID, &img.URL, &img.UserID, &img.PostID, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func (m *ImageModel) DeleteImage(imageID int, userID int) error {
	query := `
		DELETE FROM images
		WHERE id = ? AND user_id = ?
	`
	_, err := m.DB.Exec(query, imageID, userID)
	return err
}
