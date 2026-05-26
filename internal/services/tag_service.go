package services

import (
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