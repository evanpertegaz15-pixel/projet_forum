package utils

import (
    "html/template"
    "net/http"
)

func Render(w http.ResponseWriter, tmpl string, data any) {
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        http.Error(w, "Erreur de rendu.", http.StatusInternalServerError)
        return
    }
    if err := t.Execute(w, data); err != nil {
        http.Error(w, "Erreur de rendu.", http.StatusInternalServerError)
    }
}
