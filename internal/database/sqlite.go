package database

import (
    "database/sql" // API pour parler aux db
    "log" // sert à arrêter le programme si la db ne s'ouvre pas
    _ "modernc.org/sqlite"
)

func ConnectDB(path string) *sql.DB {
    db, err := sql.Open("sqlite", path) // tente d'ouvrir la db
    if err != nil {
        log.Fatalf("Erreur ouverture DB: %v", err)
    }
    _, err = db.Exec("PRAGMA foreign_keys = ON;") // activer les contraintes de foreign keys (désactivées par défaut)
    if err != nil {
        log.Fatalf("Impossible d'activer les foreign keys: %v", err)
    }
    return db
}