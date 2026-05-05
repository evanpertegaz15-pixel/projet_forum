package handlers
// users -> register, login, logout
import (
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
        SameSite: http.SameSiteLaxMode,
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
    if err != nil || user == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    tmpl := template.Must(template.ParseFiles("./internal/templates/profile.html"))
    tmpl.Execute(w, user)
}