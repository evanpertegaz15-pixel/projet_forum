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

    hashAdmin, _ := utils.HashPassword("admin")
    hashCG, _ := utils.HashPassword("clevergirl")
    hashES, _ := utils.HashPassword("estherstone")
    hashBen, _ := utils.HashPassword("ptitebosse")
    hashDarius, _ := utils.HashPassword("geekosaure")
    _, err = db.Exec(`
        INSERT OR IGNORE INTO users (email, username, password_hash, profile_picture)
        VALUES
            ('admin@dj.com', 'admin', ?, 'https://s1.qwant.com/thumbr/474x296/0/a/083d74cc5088d961e04d775a5a5a17b7cecf28102ce857b5fb0a561a9a2309/OIP.1WcjWvHMH_HLhv_YmKVX8wHaEo.jpg?u=https%3A%2F%2Ftse.mm.bing.net%2Fth%2Fid%2FOIP.1WcjWvHMH_HLhv_YmKVX8wHaEo%3Fpid%3DApi&q=0&b=1&p=0&a=0'),
            ('clever.girl@dj.com', 'PT1TEFUTEE86', ?, 'https://static.wikia.nocookie.net/jurassicpark/images/4/4a/Ronnie_1.png/revision/latest?cb=20250707082905&path-prefix=fr'),
            ('esther.stone@dj.com', 'Esth3rSt0nE', ?, 'https://static.wikia.nocookie.net/jurassicpark/images/3/38/Brooklyn%28jw%29.png/revision/latest?cb=20211024165607&path-prefix=fr'),
            ('ben.pincus@dj.com', 'BoBen', ?, 'https://static.wikia.nocookie.net/jurassicpark/images/6/6c/Bumpy_and_Ben_render.webp/revision/latest?cb=20250419043934&path-prefix=fr'),
            ('geek.osaure@dj.com', 'G3ek05aure', ?, 'https://static.wikia.nocookie.net/jurassicpark/images/8/85/DariusAdultChaosTheory.webp/revision/latest?cb=20240329090825&path-prefix=fr')
            `, hashAdmin, hashCG, hashES, hashBen, hashDarius)
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
            ('banner_url', '/static/banner.png'),
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