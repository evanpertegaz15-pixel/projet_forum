package utils

import "net/http"

func ErrorPage(w http.ResponseWriter, status int, message string) {
	titles := map[int]string{
	http.StatusBadRequest:         "Requête incorrecte",
	http.StatusUnauthorized:       "Non authentifié",
	http.StatusForbidden:          "Accès refusé",
	http.StatusMethodNotAllowed:   "Méthode non autorisée",
	http.StatusInternalServerError: "Erreur interne",
}

	title, exists := titles[status]
	if !exists {
		title = "Erreur"
	}

	w.WriteHeader(status)

	Render(w, "./internal/templates/error.html", map[string]any{
		"Code":    status,
		"Title":   title,
		"Message": message,
	})
}

func ErrorBadRequest(w http.ResponseWriter, message string) {
	ErrorPage(w, http.StatusBadRequest, message)
}

func ErrorUnauthorized(w http.ResponseWriter, message string) {
	ErrorPage(w, http.StatusUnauthorized, message)
}

func ErrorForbidden(w http.ResponseWriter, message string) {
	ErrorPage(w, http.StatusForbidden, message)
}

func ErrorMethodNotAllowed(w http.ResponseWriter, message string) {
	ErrorPage(w, http.StatusMethodNotAllowed, message)
}

func ErrorNotFound(w http.ResponseWriter, message string) {
	ErrorPage(w, http.StatusNotFound, message)
}

func ErrorInternal(w http.ResponseWriter, message string) {
	ErrorPage(w, http.StatusInternalServerError, message)
}