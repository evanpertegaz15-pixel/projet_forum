package utils

import (
    "net/http"
    "time"
)

func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
    http.SetCookie(w, &http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        MaxAge:   maxAge,
    })
}

func GetCookie(r *http.Request, name string) string {
    c, err := r.Cookie(name)
    if err != nil {
        return ""
    }
    return c.Value
}

func DeleteCookie(w http.ResponseWriter, name string) {
    http.SetCookie(w, &http.Cookie{
        Name:     name,
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })
}
