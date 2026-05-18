package handlers

import (
	"net/http"
	"strconv"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/services"
	"forum-dark-jurassic/internal/utils"
)

type UserManagementHandler struct {
	UserMgmt *services.UserManagementService
	Auth *services.AuthService
}

func NewUserManagementHandler(userMgmt *services.UserManagementService, auth *services.AuthService) *UserManagementHandler {
	return &UserManagementHandler{
		UserMgmt: userMgmt,
		Auth: auth,
	}
}

func (handler *UserManagementHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	admin, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	if !admin.HasRole("admin") {
		utils.ErrorForbidden(w, "Accès refusé")
		return
	}
	users, err := handler.UserMgmt.GetAllUsers()
	if err != nil {
		utils.ErrorInternal(w, "Erreur lors du chargement des utilisateurs")
		return
	}
	roles, err := handler.UserMgmt.GetAllRoles()
	if err != nil {
		utils.ErrorInternal(w, "Erreur lors du chargement des rôles")
		return
	}
	data := struct {
		Users []models.User
		Roles []models.Role
	}{
		Users: users,
		Roles: roles,
	}
	utils.Render(w, "./internal/templates/users_management.html", data)
}

func (handler *UserManagementHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	admin, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	if !admin.HasRole("admin") {
		utils.ErrorForbidden(w, "Accès refusé")
		return
	}
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée")
		return
	}
	err := r.ParseForm()
	if err != nil {
		utils.ErrorBadRequest(w, "Formulaire incorrect")
		return
	}
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID utilisateur invalide")
		return
	}
	roleID, err := strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID rôle invalide")
		return
	}
	err = handler.UserMgmt.AssignRole(userID, roleID)
	if err != nil {
		utils.ErrorInternal(w, "Erreur lors de l'attribution du rôle")
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (handler *UserManagementHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	admin, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	if !admin.HasRole("admin") {
		utils.ErrorForbidden(w, "Accès refusé")
		return
	}
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée")
		return
	}
	err := r.ParseForm()
	if err != nil {
		utils.ErrorBadRequest(w, "Formulaire incorrect")
		return
	}
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID utilisateur invalide")
		return
	}
	roleID, err := strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID rôle invalide")
		return
	}
	err = handler.UserMgmt.RemoveRole(userID, roleID)
	if err != nil {
		utils.ErrorInternal(w, "Erreur lors du retrait du rôle")
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (handler *UserManagementHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	admin, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	if !admin.HasRole("admin") {
		utils.ErrorForbidden(w, "Accès refusé")
		return
	}
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée")
		return
	}
	err := r.ParseForm()
	if err != nil {
		utils.ErrorBadRequest(w, "Formulaire incorrect")
		return
	}
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID utilisateur invalide")
		return
	}
	err = handler.UserMgmt.BlockUser(userID)
	if err != nil {
		utils.ErrorInternal(w, "Erreur lors du blocage de l'utilisateur")
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (handler *UserManagementHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	admin, ok := RequireAuth(w, r, handler.Auth)
	if !ok {
		return
	}
	if !admin.HasRole("admin") {
		utils.ErrorForbidden(w, "Accès refusé")
		return
	}
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée")
		return
	}
	err := r.ParseForm()
	if err != nil {
		utils.ErrorBadRequest(w, "Formulaire incorrect")
		return
	}
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.ErrorBadRequest(w, "ID utilisateur invalide")
		return
	}
	err = handler.UserMgmt.UnblockUser(userID)
	if err != nil {
		utils.ErrorInternal(w, "Erreur lors du déblocage de l'utilisateur")
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
