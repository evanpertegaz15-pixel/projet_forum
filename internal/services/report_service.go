package services

import (
	"errors"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/utils"
)

type ReportService struct {
	Reports *models.ReportModel
	Posts   *models.PostModel
	Topics  *models.TopicModel
	Users   *models.UserModel
}

func NewReportService(
	reports *models.ReportModel,
	posts *models.PostModel,
	topics *models.TopicModel,
	users *models.UserModel,
) *ReportService {
	return &ReportService{
		Reports: reports,
		Posts:   posts,
		Topics:  topics,
		Users:   users,
	}
}

func (s *ReportService) CreateReport(user *models.User, targetType string, targetID int, reason string) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}
	if targetType != "post" && targetType != "topic" {
		return errors.New("type de cible invalide")
	}
	if targetType == "post" {
		post, err := s.Posts.GetPostByID(targetID)
		if err != nil {
			return err
		}
		if post.ID == 0 {
			return errors.New("post introuvable")
		}
	} else if targetType == "topic" {
		topic, err := s.Topics.GetTopicByID(targetID)
		if err != nil {
			return err
		}
		if topic.ID == 0 {
			return errors.New("topic introuvable")
		}
	}
	return s.Reports.CreateReport(user.ID, targetType, targetID, reason)
}

func (s *ReportService) GetAllReports(user *models.User) ([]models.Report, error) {
	if !s.canManageReports(user) {
		return nil, errors.New("permission refusée")
	}
	reports, err := s.Reports.GetAllReports()
	if err != nil {
		return nil, err
	}
	for i := range reports {
		reports[i].CreatedAtAgo = utils.TimeAgo(reports[i].CreatedAt)
	}
	return reports, nil
}

func (s *ReportService) GetOpenReports(user *models.User) ([]models.Report, error) {
	if !s.canManageReports(user) {
		return nil, errors.New("permission refusée")
	}
	reports, err := s.Reports.GetOpenReports()
	if err != nil {
		return nil, err
	}
	for i := range reports {
		reports[i].CreatedAtAgo = utils.TimeAgo(reports[i].CreatedAt)
	}
	return reports, nil
}

func (s *ReportService) ResolveReport(user *models.User, reportID int) error {
	if !s.canManageReports(user) {
		return errors.New("permission refusée")
	}
	return s.Reports.ResolveReport(reportID)
}

func (s *ReportService) DeleteReport(user *models.User, reportID int) error {
	if !user.HasRole("admin") {
		return errors.New("permission refusée")
	}
	return s.Reports.DeleteReport(reportID)
}

func (s *ReportService) DeleteReportedContent(user *models.User, reportID int) error {
	if !s.canManageReports(user) {
		return errors.New("permission refusée")
	}
	report, err := s.Reports.GetReportByID(reportID)
	if err != nil {
		return err
	}
	if report == nil {
		return errors.New("report introuvable")
	}
	if report.TargetType == "post" {
		return s.Posts.Delete(report.TargetID)
	} else if report.TargetType == "topic" {
		return s.Topics.DeleteTopic(report.TargetID)
	}
	return errors.New("type de cible non supporté")
}

func (s *ReportService) canManageReports(user *models.User) bool {
	return user != nil && (user.HasRole("admin") || user.HasRole("moderator"))
}