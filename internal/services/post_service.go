package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
)

type PostService struct {
	Posts       *models.PostModel
}

func NewPostService(posts *models.PostModel) *PostService {
	return &PostService{Posts: posts}
}

func (service *PostService) GetPostsByTopic(topicID int) ([]models.Post, error) {
    return service.Posts.GetPostsByTopic(topicID)
}

func (service *PostService) CreatePost(topicID, userID int, content string) (int, error) {
    if content == "" {
        return 0, errors.New("Le contenu ne peut pas être vide.")
    }
    return service.Posts.CreatePost(topicID, userID, content)
}

func (service *PostService) CreateReply(topicID, userID, parentID int, content string) (int, error) {
    if content == "" {
        return 0, errors.New("Le contenu ne peut pas être vide.")
    }
    return service.Posts.CreateReply(topicID, userID, parentID, content)
}

func (service *PostService) GetReplies(postID int) ([]models.Post, error) {
    return service.Posts.GetReplies(postID)
}

func (service *PostService) GetPostsWithRepliesByTopic(topicID int) ([]models.PostWithReplies, error) {
    return service.Posts.GetPostsWithRepliesByTopic(topicID)
}

/*
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
*/