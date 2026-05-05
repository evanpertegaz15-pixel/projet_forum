package services
// créer, supprimer un tag, auto-complétion, validation d'unicité, association post <-> tag
import (
	//"errors"
	//"strings"
	"forum-dark-jurassic/internal/models"
)

type TagService struct {
	Tags     *models.TagModel
	PostTags *models.PostTagModel
}

func NewTagService(tags *models.TagModel, postTags *models.PostTagModel) *TagService {
	return &TagService{
		Tags:     tags,
		PostTags: postTags,
	}
}

/*
// Créer un tag avec validation d’unicité
func (s *TagService) CreateTag(name string) (int, error) {
	name = strings.TrimSpace(strings.ToLower(name))

	if name == "" {
		return 0, errors.New("le nom du tag ne peut pas être vide")
	}

	// Vérifier si le tag existe déjà
	existing, err := s.Tags.GetTagByName(name)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("ce tag existe déjà")
	}

	return s.Tags.CreateTag(name)
}

// Supprimer un tag (admin)
func (s *TagService) DeleteTag(tagID int) error {
	if tagID <= 0 {
		return errors.New("id invalide")
	}

	// Optionnel : nettoyer les relations post_tags avant suppression
	posts, err := s.PostTags.GetPostsByTagID(tagID)
	if err != nil {
		return err
	}

	for _, postID := range posts {
		_ = s.PostTags.RemoveTagFromPost(postID, tagID)
	}

	return s.Tags.DeleteTag(tagID)
}

// Associer un tag à un post
func (s *TagService) AddTagToPost(postID int, tagName string) error {
	tagName = strings.TrimSpace(strings.ToLower(tagName))

	if tagName == "" {
		return errors.New("tag vide")
	}

	// Vérifier si le tag existe
	tag, err := s.Tags.GetTagByName(tagName)
	if err != nil {
		return err
	}

	// Si le tag n'existe pas, on le crée automatiquement
	if tag == nil {
		tagID, err := s.Tags.CreateTag(tagName)
		if err != nil {
			return err
		}
		tag = &models.Tag{
			ID:   tagID,
			Name: tagName,
		}
	}

	// Vérifier si déjà associé (évite doublons)
	postTags, err := s.PostTags.GetTagsByPostID(postID)
	if err != nil {
		return err
	}

	for _, existingTagID := range postTags {
		if existingTagID == tag.ID {
			return errors.New("tag déjà associé à ce post")
		}
	}

	return s.PostTags.AddTagToPost(postID, tag.ID)
}

// Retirer un tag d’un post
func (s *TagService) RemoveTagFromPost(postID, tagID int) error {
	if postID <= 0 || tagID <= 0 {
		return errors.New("id invalide")
	}

	return s.PostTags.RemoveTagFromPost(postID, tagID)
}

// Auto-complétion de tags (simple version)
func (s *TagService) AutocompleteTags(query string) ([]models.Tag, error) {
	query = strings.TrimSpace(strings.ToLower(query))

	if len(query) < 2 {
		return []models.Tag{}, nil
	}

	allTags, err := s.Tags.GetAllTags()
	if err != nil {
		return nil, err
	}

	var result []models.Tag

	for _, tag := range allTags {
		if strings.Contains(tag.Name, query) {
			result = append(result, tag)
		}
	}

	return result, nil
}*/