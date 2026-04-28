package handlers
// users -> register, login, logout
import (
	"encoding/json"
	"net/http"
	"forum-dark-jurassic/internal/services"
)

type AuthHandler struct {
	Auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{Auth: auth}
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulaire incorrect.", http.StatusBadRequest)
        return
    }
    email := r.FormValue("email")
    username := r.FormValue("username")
    password := r.FormValue("password")
    id, err := handler.Auth.Register(email, username, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulaire incorrect.", http.StatusBadRequest)
        return
    }
    email := r.FormValue("email")
    password := r.FormValue("password")
    sessionID, err := handler.Auth.Login(email, password)
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
    http.Redirect(w, r, "/", http.StatusSeeOther)
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
        http.Error(w, "Vous n'êtes pas connecté.", http.StatusUnauthorized)
        return
    }
    user, err := handler.Auth.GetUserFromSession(cookie.Value)
    if err != nil || user == nil {
        http.Error(w, "Session incorrecte.", http.StatusUnauthorized)
        return
    }
    w.Write([]byte("Connecté en tant que : " + user.Username)) // pour l'instant, juste un affichage en texte
}