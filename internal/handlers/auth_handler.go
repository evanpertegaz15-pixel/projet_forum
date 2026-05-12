package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"forum-dark-jurassic/internal/services"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ─── Google OAuth Config ──────────────────────────────────────────────────────

var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

type googleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ─── Handler Struct ───────────────────────────────────────────────────────────

type AuthHandler struct {
	Auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{Auth: auth}
}

// ─── Original Handlers ────────────────────────────────────────────────────────

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
		MaxAge:   7 * 24 * 60 * 60,
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
	tmpl := template.Must(template.ParseFiles("./internal/templates/profile.html"))
	tmpl.Execute(w, user)
}

// ─── Google OAuth Handlers ────────────────────────────────────────────────────

func (handler *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a random CSRF state token
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Save state in a short-lived cookie to verify on callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	// Redirect user to Google's login page
	url := googleOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (handler *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Validate state to prevent CSRF attacks
	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "État OAuth invalide.", http.StatusBadRequest)
		return
	}

	// 2. Delete the state cookie immediately so it can't be reused
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		MaxAge: -1,
		Path:   "/",
	})

	// 3. Exchange the temporary code for a real access token
	token, err := googleOAuthConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Échange de token échoué.", http.StatusInternalServerError)
		return
	}

	// 4. Use the token to fetch user info from Google
	client := googleOAuthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Impossible de récupérer les infos utilisateur.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var userInfo googleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		http.Error(w, "Erreur de décodage des infos utilisateur.", http.StatusInternalServerError)
		return
	}

	// 5. Find existing user by email or create a new one
	userID, err := handler.Auth.FindOrCreateGoogleUser(userInfo.Email, userInfo.Name)
	if err != nil {
		http.Error(w, "Erreur utilisateur.", http.StatusInternalServerError)
		return
	}

	// 6. Create a session for that user
	sessionID, err := handler.Auth.CreateSession(userID)
	if err != nil {
		http.Error(w, "Erreur de session.", http.StatusInternalServerError)
		return
	}

	// 7. Set the session cookie — same as Login
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
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
    user, err := handler.Auth.GetUserFromSession(cookie.Value)
    if err != nil {
        log.Printf("DeleteAccount: failed to get user from session: %v\n", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    if user == nil {
        log.Printf("DeleteAccount: no user found for session %s\n", cookie.Value)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    err = handler.Auth.Users.DeleteUser(user.ID)
    if err != nil {
        log.Printf("DeleteAccount: DeleteUser error user_id=%d: %v\n", user.ID, err)
        http.Error(w, "Impossible de supprimer le compte.", http.StatusInternalServerError)
        return
    }
    if err := handler.Auth.Sessions.DeleteSession(cookie.Value); err != nil {
        log.Printf("DeleteAccount: DeleteSession error session=%s user_id=%d: %v\n", cookie.Value, user.ID, err)
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
