package handlers

import (
    "net/http"
	"forum-dark-jurassic/internal/models"
    "forum-dark-jurassic/internal/services"
    "forum-dark-jurassic/internal/utils"
)

func RequireAuth(w http.ResponseWriter, r *http.Request, auth *services.AuthService) (*models.User, bool) {
	sessionID := utils.GetCookie(r, "session_id")
    if sessionID == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return nil, false
    }
    user, err := auth.GetUserFromSession(sessionID)
    if err != nil || user == nil {
		utils.DeleteCookie(w, "session_id")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return nil, false
    }
    return user, true
}

func RequireRole(auth *services.AuthService, role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := RequireAuth(w, r, auth)
            if !ok {
                return
            }
            if !auth.UserHasRole(user, role) {
                utils.ErrorForbidden(w, "Permission refusée.")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}