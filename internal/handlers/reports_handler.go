package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/services"
	"forum-dark-jurassic/internal/utils"
)

type ReportHandler struct {
	Reports *services.ReportService
	Auth    *services.AuthService
}

func NewReportHandler(reportService *services.ReportService, auth *services.AuthService) *ReportHandler {
	return &ReportHandler{
		Reports: reportService,
		Auth:    auth,
	}
}

func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
		return
	}
	user, ok := RequireAuth(w, r, h.Auth)
	if !ok {
		return
	}
	targetType := r.FormValue("target_type")
	if targetType != "post" && targetType != "topic" {
		utils.ErrorBadRequest(w, "Type de cible invalide.")
		return
	}
	targetID, err := strconv.Atoi(r.FormValue("target_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID de cible invalide.")
		return
	}
	reason := r.FormValue("reason")
	if reason == "" {
		utils.ErrorBadRequest(w, "Raison requise.")
		return
	}
	returnURL := r.FormValue("return_url")
	if returnURL == "" {
		returnURL = r.Referer()
	}
	if returnURL == "" {
		returnURL = "/"
	}
	err = h.Reports.CreateReport(user, targetType, targetID, reason)
	if err != nil {
		http.Redirect(w, r, addQueryParam(returnURL, "report_status", "error"), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, addQueryParam(returnURL, "report_status", "success"), http.StatusSeeOther)
}

func addQueryParam(rawURL, key, value string) string {
	if rawURL == "" {
		return "/"
	}
	if strings.Contains(rawURL, "?") {
		return rawURL + "&" + url.QueryEscape(key) + "=" + url.QueryEscape(value)
	}
	return rawURL + "?" + url.QueryEscape(key) + "=" + url.QueryEscape(value)
}

func (h *ReportHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	user, ok := RequireAuth(w, r, h.Auth)
	if !ok {
		return
	}
	reports, err := h.Reports.GetAllReports(user)
	if err != nil {
		utils.ErrorForbidden(w, "Permission refusée.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) GetOpenReports(w http.ResponseWriter, r *http.Request) {
	user, ok := RequireAuth(w, r, h.Auth)
	if !ok {
		return
	}
	reports, err := h.Reports.GetOpenReports(user)
	if err != nil {
		utils.ErrorForbidden(w, "Permission refusée.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) ResolveReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
		return
	}
	user, ok := RequireAuth(w, r, h.Auth)
	if !ok {
		return
	}
	reportID, err := strconv.Atoi(r.FormValue("report_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID de report invalide.")
		return
	}
	returnURL := r.FormValue("return_url")
	if returnURL == "" {
		returnURL = "/reports"
	}
	err = h.Reports.ResolveReport(user, reportID)
	if err != nil {
		http.Redirect(w, r, addQueryParam(returnURL, "report_status", "error"), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, addQueryParam(returnURL, "report_status", "success"), http.StatusSeeOther)
}

func (h *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
		return
	}
	user, ok := RequireAuth(w, r, h.Auth)
	if !ok {
		return
	}
	reportID, err := strconv.Atoi(r.FormValue("report_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID de report invalide.")
		return
	}
	returnURL := r.FormValue("return_url")
	if returnURL == "" {
		returnURL = "/reports"
	}
	err = h.Reports.DeleteReport(user, reportID)
	if err != nil {
		http.Redirect(w, r, addQueryParam(returnURL, "report_status", "error"), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, addQueryParam(returnURL, "report_status", "success"), http.StatusSeeOther)
}

func (h *ReportHandler) DeleteReportedContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
		return
	}
	user, ok := RequireAuth(w, r, h.Auth)
	if !ok {
		return
	}
	reportID, err := strconv.Atoi(r.FormValue("report_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID de report invalide.")
		return
	}
	err = h.Reports.DeleteReportedContent(user, reportID)
	if err != nil {
		utils.ErrorForbidden(w, "Permission refusée.")
		return
	}
	returnURL := r.FormValue("return_url")
	if returnURL == "" {
		returnURL = "/reports"
	}
	http.Redirect(w, r, addQueryParam(returnURL, "report_status", "success"), http.StatusSeeOther)
}

func (h *ReportHandler) ShowReports(w http.ResponseWriter, r *http.Request) {
	user, ok := RequireAuth(w, r, h.Auth)
	if !ok {
		return
	}
	reports, err := h.Reports.GetOpenReports(user)
	if err != nil {
		utils.ErrorForbidden(w, "Permission refusée.")
		return
	}
	for i := range reports {
		reports[i].CreatedAtAgo = utils.TimeAgo(reports[i].CreatedAt)
	}
	data := struct {
		User    *models.User
		Reports []models.Report
		Status  string
	}{
		User:    user,
		Reports: reports,
		Status:  r.URL.Query().Get("report_status"),
	}
	utils.Render(w, "./internal/templates/reports.html", data)
}