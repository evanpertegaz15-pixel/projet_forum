package services

import (
	"forum-dark-jurassic/internal/models"
)

type CategoryService struct {
	Categories *models.CategoryModel
}

func NewCategoryService(categories *models.CategoryModel) *CategoryService {
	return &CategoryService{Categories: categories}
}

func (service *CategoryService) GetAllCategories() ([]models.Category, error) {
    return service.Categories.GetAllCategories()
}

func (service *CategoryService) GetCategory(id int) (*models.Category, error) {
    return service.Categories.GetCategoryByID(id)
}