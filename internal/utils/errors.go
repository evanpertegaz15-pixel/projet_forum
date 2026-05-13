package utils

import (
    "net/http"
)

func ErrorBadRequest(w http.ResponseWriter, message string) {
    http.Error(w, message, http.StatusBadRequest)
}

func ErrorUnauthorized(w http.ResponseWriter, message string) {
    http.Error(w, message, http.StatusUnauthorized)
}

func ErrorForbidden(w http.ResponseWriter, message string) {
    http.Error(w, message, http.StatusForbidden)
}

func ErrorMethodNotAllowed(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusMethodNotAllowed)
}

func ErrorNotFound(w http.ResponseWriter, message string) {
    http.Error(w, message, http.StatusNotFound)
}

func ErrorInternal(w http.ResponseWriter, message string) {
    http.Error(w, message, http.StatusInternalServerError)
}
