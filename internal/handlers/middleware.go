package handlers

import (
    "net/http"
	"forum-dark-jurassic/internal/models"
    "forum-dark-jurassic/internal/services"
)

func RequireAuth(w http.ResponseWriter, r *http.Request, auth *services.AuthService) (*models.User, bool) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return nil, false
    }

    user, err := auth.GetUserFromSession(cookie.Value)
    if err != nil || user == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return nil, false
    }

    return user, true
}

/*
type Middleware struct {
	Auth *services.AuthService
}

func NewMiddleware(authService *services.AuthService) *Middleware {
	return &Middleware{
		Auth: authService,
	}
}

// LoadUserFromSession : injecte userID dans le context si session valide
func (m *Middleware) LoadUserFromSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.Auth.GetUserFromSession(cookie.Value)
		if err != nil || user == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", user.ID)
		ctx = context.WithValue(ctx, "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole : vérifie un rôle précis (admin, mod, etc.)
func (m *Middleware) RequireRole(roleName string, roleService RoleChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID")

			if userID == nil {
				http.Error(w, "Non authentifié", http.StatusUnauthorized)
				return
			}

			ok, err := roleService.UserHasRoleByName(userID.(int), roleName)
			if err != nil {
				http.Error(w, "Erreur serveur", http.StatusInternalServerError)
				return
			}

			if !ok {
				http.Error(w, "Permission refusée", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Interface pour découpler RoleModel / RoleService
type RoleChecker interface {
	UserHasRoleByName(userID int, roleName string) (bool, error)
}*/