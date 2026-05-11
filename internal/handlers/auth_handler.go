package handlers

import (
    "log"
	"html/template"
	"net/http"
	"forum-dark-jurassic/internal/services"
)

type AuthHandler struct {
	Auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{Auth: auth}
}

func Home(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("./internal/templates/home.html"))
    tmpl.Execute(w, nil)
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulaire incorrect.", http.StatusBadRequest)
        return
    }
    email := r.FormValue("email")
    username := r.FormValue("username")
    password := r.FormValue("password")
    _, err := handler.Auth.Register(email, username, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("./internal/templates/register.html"))
    tmpl.Execute(w, nil)
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulaire incorrect.", http.StatusBadRequest)
        return
    }
    identifier := r.FormValue("identifier")
    password := r.FormValue("password")
    sessionID, err := handler.Auth.Login(identifier, password)
    if err != nil {
        http.Error(w, "Identifiants incorrects.", http.StatusUnauthorized)
        return
    }
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        MaxAge:   7 * 24 * 60 * 60, // 7 days
    })
    http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("./internal/templates/login.html"))
    tmpl.Execute(w, nil)
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == nil {
        _ = handler.Auth.Logout(cookie.Value)
    }
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        HttpOnly: true,
        SameSite: http.SameSiteDefaultMode,
    })
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    user, err := handler.Auth.GetUserFromSession(cookie.Value)
    if err != nil {
        log.Printf("Error getting user from session: %v\n", err)
    }
    if err != nil || user == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    log.Printf("COOKIE:", cookie.Value)
    tmpl := template.Must(template.ParseFiles("./internal/templates/profile.html"))
    tmpl.Execute(w, user)
}

func (handler *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée.", http.StatusMethodNotAllowed)
        return
    }
    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    user, err := h.Auth.GetUserFromSession(cookie.Value)
    if err != nil || user == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    err = handler.Auth.Users.DeleteUser(user.ID)
    if err != nil {
        http.Error(w, "Impossible de supprimer le compte.", http.StatusInternalServerError)
        return
    }
    h.Auth.Sessions.DeleteSession(cookie.Value)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}