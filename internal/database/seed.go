package database

import (
    "database/sql"
    "log"
    "forum-dark-jurassic/internal/utils"
)

func Seed(db *sql.DB) {
    _, err := db.Exec(`INSERT OR IGNORE INTO roles (name, label) VALUES ('admin', 'Administrateur'), ('moderator', 'Modérateur'), ('ranger', 'Ranger'), ('user', 'Utilisateur')`)
    if err != nil {
        log.Println("Erreur seed rôles :", err)
    }

    hash, _ := utils.HashPassword("admin")
    _, err = db.Exec(`
        INSERT OR IGNORE INTO users (email, username, password_hash)
        VALUES ('admin@dj.com', 'admin', ?)`, hash)
    if err != nil {
        log.Println("Erreur seed admin :", err)
    }

    _, err = db.Exec(`
        INSERT OR IGNORE INTO categories (name) VALUES ('Général'), ('Isla Nublar'), ('Herbivores'), ('Omnivores'), ('Carnivores'), ('Théories'), ('Incidents')`)
    if err != nil {
        log.Println("Erreur seed categories :", err)
    }

    _, err = db.Exec(`
        INSERT OR IGNORE INTO settings (key, value) VALUES
            ('forum_name', 'Dark Jurassic'),
            ('forum_description', 'Le forum des théories les plus folles en rapport avec les dinosaures.'),
            ('banner_url', '/internal/templates/static/banner.png'),
            ('admin_email', 'admin@dj.com'),
            ('maintenance_mode', 'false'),
            ('allow_registration', 'true'),
            ('theme', 'dark'),
            ('max_upload_size', '20971520'), -- 20 MB
            ('footer_text', '© 2026 Dark Jurassic'),
            ('default_role', 'user'),
            ('allow_post_images', 'true'),
            ('allow_post_edit', 'true')
        `)
    if err != nil {
        log.Println("Erreur seed settings :", err)
    }
}