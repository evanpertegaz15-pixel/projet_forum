package utils

import (
    "errors"
    "regexp"
    "strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
func ValidateEmail(email string) bool {
    email = strings.TrimSpace(email)
    return emailRegex.MatchString(email)
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
func ValidateUsername(username string) bool {
    username = strings.TrimSpace(username)
    return usernameRegex.MatchString(username)
}

func ValidatePassword(password string) error {
    password = strings.TrimSpace(password)
    if len(password) < 6 {
        return errors.New("Le mot de passe doit faire au moins 6 caractères.")
    }
    if len(password) > 72 {
        return errors.New("Le mot de passe est trop long.")
    }
    return nil
}

func ValidatePostTitle(title string) error {
    title = strings.TrimSpace(title)
    if len(title) == 0 {
        return errors.New("Le titre ne peut pas être vide.")
    }
    if len(title) < 3 {
        return errors.New("Le titre doit faire au moins 3 caractères.")
    }
    if len(title) > 120 {
        return errors.New("Le titre est trop long.")
    }
    return nil
}