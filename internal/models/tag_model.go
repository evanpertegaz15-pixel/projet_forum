package models
// lister les tags, les associer à un poste, les créer en admin
import (
	"database/sql"
)

type Tag struct {
	ID   int
	Name string
}

type TagModel struct {
	DB *sql.DB
}

func NewTagModel(db *sql.DB) *TagModel {
	return &TagModel{DB: db}
}

// Créer un tag (admin)
func (m *TagModel) CreateTag(name string) (int, error) {
	query := `
		INSERT INTO tags (name)
		VALUES (?)
	`

	result, err := m.DB.Exec(query, name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Lister tous les tags
func (m *TagModel) GetAllTags() ([]Tag, error) {
	query := `
		SELECT id, name
		FROM tags
		ORDER BY name ASC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

// Récupérer un tag par ID
func (m *TagModel) GetTagByID(tagID int) (*Tag, error) {
	query := `
		SELECT id, name
		FROM tags
		WHERE id = ?
	`

	var t Tag

	err := m.DB.QueryRow(query, tagID).Scan(&t.ID, &t.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Récupérer un tag par nom (utile pour auto-tagging / recherche)
func (m *TagModel) GetTagByName(name string) (*Tag, error) {
	query := `
		SELECT id, name
		FROM tags
		WHERE name = ?
	`

	var t Tag

	err := m.DB.QueryRow(query, name).Scan(&t.ID, &t.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Supprimer un tag (admin)
func (m *TagModel) DeleteTag(tagID int) error {
	query := `
		DELETE FROM tags
		WHERE id = ?
	`

	_, err := m.DB.Exec(query, tagID)
	return err
}