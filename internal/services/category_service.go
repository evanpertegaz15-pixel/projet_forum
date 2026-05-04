package services
// un admin peut créer, supprimer, renommer une catégorie, validation d'unicité, association post <-> catégorie

package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
)

type CategoryService struct {
	Categories *models.CategoryModel
	Posts      *models.PostModel
}

func NewCategoryService(categories *models.CategoryModel, posts *models.PostModel) *CategoryService {
	return &CategoryService{
		Categories: categories,
		Posts:      posts,
	}
}

// Vérification admin (simple)
func isAdmin(user *models.User) bool {
	return user != nil && user.Role == "admin"
}

// Créer une catégorie
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

// Renommer une catégorie
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

// Supprimer une catégorie
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
}

// Associer un post à une catégorie
func (s *CategoryService) AddPostToCategory(postID, categoryID int) error {
	return s.Posts.AddCategory(postID, categoryID)
}

// Retirer un post d’une catégorie
func (s *CategoryService) RemovePostFromCategory(postID, categoryID int) error {
	return s.Posts.RemoveCategory(postID, categoryID)
}

// Lister les catégories
func (s *CategoryService) GetAllCategories() ([]*models.Category, error) {
	return s.Categories.GetAll()
}

// Obtenir les posts d’une catégorie
func (s *CategoryService) GetPostsByCategory(categoryID int) ([]*models.Post, error) {
	return s.Posts.GetByCategory(categoryID)
}