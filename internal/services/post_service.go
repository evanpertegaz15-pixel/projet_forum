// création de post, édition, suppression, permissions, tags, images, catégories

package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
)

type PostService struct {
	Posts       *models.PostModel
	Categories  *models.CategoryModel
	Tags        *models.TagModel
	Images      *ImageService
}

func NewPostService(
	posts *models.PostModel,
	categories *models.CategoryModel,
	tags *models.TagModel,
	images *ImageService,
) *PostService {
	return &PostService{
		Posts:      posts,
		Categories: categories,
		Tags:       tags,
		Images:     images,
	}
}

// admin check
func isAdminCheck(user *models.User) bool {
	return user != nil && user.Role == "admin"
}

// Créer un post
func (s *PostService) CreatePost(
	user *models.User,
	title string,
	content string,
	categoryIDs []int,
	tagNames []string,
	imageFilename *string,
) (int, error) {

	if user == nil {
		return 0, errors.New("utilisateur non connecté")
	}

	if title == "" || content == "" {
		return 0, errors.New("titre ou contenu vide")
	}

	// créer post
	postID, err := s.Posts.Create(user.ID, title, content)
	if err != nil {
		return 0, err
	}

	// catégories
	for _, catID := range categoryIDs {
		_ = s.Posts.AddCategory(postID, catID)
	}

	// tags
	for _, tag := range tagNames {
		tagID, _ := s.Tags.GetOrCreate(tag)
		_ = s.Posts.AddTag(postID, tagID)
	}

	// image
	if imageFilename != nil {
		_ = s.Posts.SetImage(postID, *imageFilename)
	}

	return postID, nil
}

// Modifier un post
func (s *PostService) UpdatePost(
	user *models.User,
	postID int,
	title string,
	content string,
	categoryIDs []int,
	tagNames []string,
	newImage *string,
) error {

	post, err := s.Posts.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post introuvable")
	}

	// permission
	if post.UserID != user.ID && !isAdminCheck(user) {
		return errors.New("permission refusée")
	}

	// update contenu
	err = s.Posts.Update(postID, title, content)
	if err != nil {
		return err
	}

	// catégories (reset)
	_ = s.Posts.ClearCategories(postID)
	for _, catID := range categoryIDs {
		_ = s.Posts.AddCategory(postID, catID)
	}

	// tags (reset)
	_ = s.Posts.ClearTags(postID)
	for _, tag := range tagNames {
		tagID, _ := s.Tags.GetOrCreate(tag)
		_ = s.Posts.AddTag(postID, tagID)
	}

	// image (remplacement)
	if newImage != nil {
		if post.Image != "" {
			_ = s.Images.DeleteImage(post.Image)
		}
		_ = s.Posts.SetImage(postID, *newImage)
	}

	return nil
}

// Supprimer un post
func (s *PostService) DeletePost(user *models.User, postID int) error {

	post, err := s.Posts.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post introuvable")
	}

	// permission
	if post.UserID != user.ID && !isAdminCheck(user) {
		return errors.New("permission refusée")
	}

	// supprimer image
	if post.Image != "" {
		_ = s.Images.DeleteImage(post.Image)
	}

	// nettoyer relations
	_ = s.Posts.ClearCategories(postID)
	_ = s.Posts.ClearTags(postID)

	return s.Posts.Delete(postID)
}

// Obtenir un post
func (s *PostService) GetPostByID(postID int) (*models.PostModel, error) {
	return s.Posts.FindByID(postID)
}

// Lister posts
func (s *PostService) GetAllPosts() ([]*models.PostModel, error) {
	return s.Posts.GetAll()
}

// Par catégorie
func (s *PostService) GetPostsByCategory(categoryID int) ([]*models.PostModel, error) {
	return s.Posts.GetByCategory(categoryID)
}

// Par tag
func (s *PostService) GetPostsByTag(tagID int) ([]*models.PostModel, error) {
	return s.Posts.GetByTag(tagID)
}
