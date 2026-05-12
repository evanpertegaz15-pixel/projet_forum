package models

import (
	"database/sql"
)

type Category struct {
	ID          int
	Name        string
	Description sql.NullString
	CreatedAt   sql.NullTime
}

type CategoryModel struct {
	DB *sql.DB
}

func NewCategoryModel(db *sql.DB) *CategoryModel {
	return &CategoryModel{DB:db}
}

func (model *CategoryModel) GetAllCategories() ([]Category, error) {
	rows, err := model.DB.Query(`SELECT id, name, description, created_at FROM categories ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (model *CategoryModel) GetCategoryByID(id int) (*Category, error) {
    row := model.DB.QueryRow(`SELECT id, name, description, created_at FROM categories WHERE id = ?`, id)
    var category Category
    err := row.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &category, nil
}