package handlers
// reports -> signaler un post / commentaire
import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum-dark-jurassic/internal/services"
)

type ReportHandler struct {
	Reports *services.ReportService
}

func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{
		Reports: reportService,
	}
}

// Signaler un post ou un commentaire
func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	targetType := r.FormValue("target_type") // "post" ou "comment"
	if targetType != "post" && targetType != "comment" {
		http.Error(w, "target_type invalide", http.StatusBadRequest)
		return
	}

	targetID, err := strconv.Atoi(r.FormValue("target_id"))
	if err != nil {
		http.Error(w, "target_id invalide", http.StatusBadRequest)
		return
	}

	reason := r.FormValue("reason")
	if reason == "" {
		http.Error(w, "raison requise", http.StatusBadRequest)
		return
	}

	err = h.Reports.CreateReport(
		userID.(int),
		targetType,
		targetID,
		reason,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Signalement envoyé",
	})
}

// Lister tous les reports (admin)
func (h *ReportHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.Reports.GetAllReports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// Lister uniquement les reports non traités (admin)
func (h *ReportHandler) GetOpenReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.Reports.GetOpenReports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// Marquer un report comme traité (admin)
func (h *ReportHandler) ResolveReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	reportID, err := strconv.Atoi(r.FormValue("report_id"))
	if err != nil {
		http.Error(w, "report_id invalide", http.StatusBadRequest)
		return
	}

	err = h.Reports.ResolveReport(reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Report traité",
	})
}

// Supprimer un report (admin nettoyage)
func (h *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	reportID, err := strconv.Atoi(r.FormValue("report_id"))
	if err != nil {
		http.Error(w, "report_id invalide", http.StatusBadRequest)
		return
	}

	err = h.Reports.DeleteReport(reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Report supprimé",
	})
}