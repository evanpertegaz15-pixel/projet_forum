// signaler un post / commentaire, modération

package services

import (
	//"errors"
	"forum-dark-jurassic/internal/models"
)

type ReportService struct {
	Reports  *models.ReportModel
	Posts    *models.PostModel
	Comments *models.CommentModel
}

func NewReportService(
	reports *models.ReportModel,
	posts *models.PostModel,
	comments *models.CommentModel,
) *ReportService {
	return &ReportService{
		Reports:  reports,
		Posts:    posts,
		Comments: comments,
	}
}

/*
// admin check
func isAdminCheck(user *models.User) bool {
	return user != nil && user.Role == "admin"
}

// Signaler un post
func (s *ReportService) ReportPost(user *models.User, postID int, reason string) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	post, err := s.Posts.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post introuvable")
	}

	// éviter spam (1 report par user par post)
	exists, err := s.Reports.Exists(user.ID, postID, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("déjà signalé")
	}

	return s.Reports.Create(user.ID, postID, 0, reason)
}

// Signaler un commentaire
func (s *ReportService) ReportComment(user *models.User, commentID int, reason string) error {
	if user == nil {
		return errors.New("utilisateur non connecté")
	}

	comment, err := s.Comments.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("commentaire introuvable")
	}

	exists, err := s.Reports.Exists(user.ID, 0, commentID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("déjà signalé")
	}

	return s.Reports.Create(user.ID, 0, commentID, reason)
}

// Lister tous les reports (admin)
func (s *ReportService) GetAllReports(user *models.User) ([]*models.ReportModel, error) {
	if !isAdminCheck(user) {
		return nil, errors.New("permission refusée")
	}
	return s.Reports.GetAll()
}

// Reports en attente
func (s *ReportService) GetPendingReports(user *models.User) ([]*models.ReportModel, error) {
	if !isAdminCheck(user) {
		return nil, errors.New("permission refusée")
	}
	return s.Reports.GetByStatus("pending")
}

// Marquer comme traité
func (s *ReportService) ResolveReport(user *models.User, reportID int) error {
	if !isAdminCheck(user) {
		return errors.New("permission refusée")
	}
	return s.Reports.UpdateStatus(reportID, "resolved")
}

// Rejeter un report
func (s *ReportService) RejectReport(user *models.User, reportID int) error {
	if !isAdminCheck(user) {
		return errors.New("permission refusée")
	}
	return s.Reports.UpdateStatus(reportID, "rejected")
}

// Action modération : supprimer contenu signalé
func (s *ReportService) DeleteReportedContent(user *models.User, reportID int) error {
	if !isAdminCheck(user) {
		return errors.New("permission refusée")
	}

	report, err := s.Reports.FindByID(reportID)
	if err != nil {
		return err
	}
	if report == nil {
		return errors.New("report introuvable")
	}

	// si post
	if report.PostID != nil {
		return s.Posts.Delete(*report.PostID)
	}

	// si commentaire
	if report.CommentID != nil {
		return s.Comments.Delete(*report.CommentID)
	}

	return nil
}
*/