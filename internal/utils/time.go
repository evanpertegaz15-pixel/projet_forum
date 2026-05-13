package utils

import (
    "time"
    "fmt"
)

func TimeAgo(t time.Time) string {
    now := time.Now()
    diff := now.Sub(t)
    seconds := int(diff.Seconds())
    minutes := seconds / 60
    hours := minutes / 60
    days := hours / 24
    switch {
    case seconds < 60:
        return "À l’instant"
    case minutes < 60:
        if minutes == 1 {
            return "Il y a 1 minute"
        }
        return fmt.Sprintf("Il y a %d minutes", minutes)
    case hours < 24:
        if hours == 1 {
            return "Il y a 1 heure"
        }
        return fmt.Sprintf("Il y a %d heures", hours)
    case days == 1:
        return "Hier"
    case days < 30:
        return fmt.Sprintf("Il y a %d jours", days)
    case days < 365:
        months := days / 30
        if months == 1 {
            return "Il y a 1 mois"
        }
        return fmt.Sprintf("Il y a %d mois", months)
    default:
        years := days / 365
        if years == 1 {
            return "Il y a 1 an"
        }
        return fmt.Sprintf("Il y a %d ans", years)
    }
}
