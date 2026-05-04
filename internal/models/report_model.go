package models
// panneau admin, enregister un report
import (
	"database/sql"
	"time"
)

type Report struct {
	ID          int
	ReporterID  int
	TargetType  string // "post", "comment", "user"
	TargetID    int
	Reason      string
	IsResolved  bool
	CreatedAt   time.Time
}

type ReportModel struct {
	DB *sql.DB
}

func NewReportModel(db *sql.DB) *ReportModel {
	return &ReportModel{DB: db}
}

// Créer un report (signalement)
func (m *ReportModel) CreateReport(reporterID int, targetType string, targetID int, reason string) error {
	query := `
		INSERT INTO reports (reporter_id, target_type, target_id, reason, is_resolved, created_at)
		VALUES (?, ?, ?, ?, 0, ?)
	`

	_, err := m.DB.Exec(query, reporterID, targetType, targetID, reason, time.Now())
	return err
}

// Lister tous les reports (admin)
func (m *ReportModel) GetAllReports() ([]Report, error) {
	query := `
		SELECT id, reporter_id, target_type, target_id, reason, is_resolved, created_at
		FROM reports
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(
			&r.ID,
			&r.ReporterID,
			&r.TargetType,
			&r.TargetID,
			&r.Reason,
			&r.IsResolved,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	return reports, nil
}

// Lister uniquement les reports non résolus
func (m *ReportModel) GetOpenReports() ([]Report, error) {
	query := `
		SELECT id, reporter_id, target_type, target_id, reason, is_resolved, created_at
		FROM reports
		WHERE is_resolved = 0
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(
			&r.ID,
			&r.ReporterID,
			&r.TargetType,
			&r.TargetID,
			&r.Reason,
			&r.IsResolved,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	return reports, nil
}

// Marquer un report comme résolu
func (m *ReportModel) ResolveReport(reportID int) error {
	query := `
		UPDATE reports
		SET is_resolved = 1
		WHERE id = ?
	`

	_, err := m.DB.Exec(query, reportID)
	return err
}

// Supprimer un report (nettoyage admin)
func (m *ReportModel) DeleteReport(reportID int) error {
	query := `
		DELETE FROM reports
		WHERE id = ?
	`

	_, err := m.DB.Exec(query, reportID)
	return err
}