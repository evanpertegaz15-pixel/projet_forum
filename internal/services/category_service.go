package services

import (
	//"errors"
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

/*
func (s *CategoryService) CreateCategory(user *models.User, name string) (int, error) {
	if !isAdmin(user) {
		return 0, errors.New("permission refusée")
	}

	existing, err := s.Categories.FindByName(name)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("cette catégorie existe déjà")
	}

	return s.Categories.Create(name)
}

func (s *CategoryService) RenameCategory(user *models.User, categoryID int, newName string) error {
	if !isAdmin(user) {
		return errors.New("permission refusée")
	}

	existing, err := s.Categories.FindByName(newName)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("nom déjà utilisé")
	}

	return s.Categories.UpdateName(categoryID, newName)
}

func (s *CategoryService) DeleteCategory(user *models.User, categoryID int) error {
	if !isAdmin(user) {
		return errors.New("permission refusée")
	}

	// option : retirer les liens avec les posts
	err := s.Posts.RemoveCategoryFromPosts(categoryID)
	if err != nil {
		return err
	}

	return s.Categories.Delete(categoryID)
}*/