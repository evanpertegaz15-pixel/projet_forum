package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"forum-dark-jurassic/internal/services"
	"forum-dark-jurassic/internal/utils"
	"io"
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
	utils.Render(w, "./internal/templates/home.html", nil)
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire incorrect.", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	if !utils.ValidateEmail(email) {
		utils.Render(w, "./internal/templates/register.html", map[string]any{
			"Error": "Email invalide.",
		})
		return
	}
	if !utils.ValidateUsername(username) {
		utils.Render(w, "./internal/templates/register.html", map[string]any{
			"Error": "Nom d'utilisateur invalide.",
		})
		return
	}
	if err := utils.ValidatePassword(password); err != nil {
		utils.Render(w, "./internal/templates/register.html", map[string]any{
			"Error": err.Error(),
		})
		return
	}
	_, err := handler.Auth.Register(email, username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (handler *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "./internal/templates/register.html", nil)
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.ErrorBadRequest(w, "Formulaire incorrect.")
		return
	}
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")
	if !utils.ValidateEmail(identifier) && !utils.ValidateUsername(identifier) {
		utils.Render(w, "./internal/templates/login.html", map[string]any{
			"Error": "Identifiant invalide.",
		})
		return
	}
	if err := utils.ValidatePassword(password); err != nil {
		utils.Render(w, "./internal/templates/login.html", map[string]any{
			"Error": err.Error(),
		})
		return
	}
	sessionID, err := handler.Auth.Login(identifier, password)
	if err != nil {
		utils.Render(w, "./internal/templates/login.html", map[string]any{
            "Error": "Identifiants incorrects.",
        })
		return
	}
	utils.SetCookie(w, "session_id", sessionID, 7*24*60*60)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (handler *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, "./internal/templates/login.html", nil)
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_ = handler.Auth.Logout(cookie.Value)
	}
	utils.DeleteCookie(w, "session_id")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := RequireAuth(w, r, handler.Auth)
	if !ok { return }
	user.CreatedAtAgo = utils.TimeAgo(user.CreatedAt)
	utils.Render(w, "./internal/templates/profile.html", user)
}

func (handler *AuthHandler) ShowEditProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := RequireAuth(w, r, handler.Auth)
	if !ok { return }
	utils.Render(w, "./internal/templates/edit_profile.html", user)
}

func (handler *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMethodNotAllowed(w, "Méthode non autorisée")
		return
	}
	user, ok := RequireAuth(w, r, handler.Auth)
	if !ok { return }
	err := r.ParseForm()
	if err != nil {
		utils.ErrorBadRequest(w, "Formulaire incorrect")
		return
	}
	newEmail := r.FormValue("email")
	newUsername := r.FormValue("username")
	newPassword := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	if newEmail != "" && newEmail != user.Email {
		if !utils.ValidateEmail(newEmail) {
			utils.ErrorBadRequest(w, "Email invalide")
			return
		}
		if err := handler.Auth.UpdateEmail(user.ID, newEmail); err != nil {
			utils.ErrorInternal(w, "Erreur lors de la mise à jour de l'email")
			return
		}
	}
	if newUsername != "" && newUsername != user.Username {
		if !utils.ValidateUsername(newUsername) {
			utils.ErrorBadRequest(w, "Nom d'utilisateur invalide")
			return
		}
		if err := handler.Auth.UpdateUsername(user.ID, newUsername); err != nil {
			utils.ErrorInternal(w, "Erreur lors de la mise à jour du nom d'utilisateur")
			return
		}
	}
	if newPassword != "" {
		if err := utils.ValidatePassword(newPassword); err != nil {
			utils.ErrorBadRequest(w, err.Error())
			return
		}
		if newPassword != confirmPassword {
			utils.ErrorBadRequest(w, "Les mots de passe ne correspondent pas")
			return
		}
		if err := handler.Auth.UpdatePassword(user.ID, newPassword); err != nil {
			utils.ErrorInternal(w, "Erreur lors de la mise à jour du mot de passe")
			return
		}
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
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
	sessionID, err := handler.Auth.CreateSessionFromGoogle(userID)
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
        utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
        return
    }
	user, ok := RequireAuth(w, r, handler.Auth)
	if !ok { return }
    _, err := r.Cookie("session_id")
    if err != nil {
        utils.ErrorUnauthorized(w, "Session introuvable.")
        return
    }
    if err := handler.Auth.Users.DeleteUser(user.ID); err != nil {
    	utils.ErrorInternal(w, "Impossible de supprimer le compte.")
        return
    }
	sessionID := utils.GetCookie(r, "session_id")
	if sessionID != "" {
		_ = handler.Auth.Sessions.DeleteSession(sessionID)
	}
	utils.DeleteCookie(w, "session_id")
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
