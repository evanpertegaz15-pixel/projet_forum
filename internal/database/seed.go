package database

import (
    "database/sql"
    "log"
    "forum-dark-jurassic/internal/utils"
)

func Seed(db *sql.DB) {
    _, err := db.Exec(`INSERT OR IGNORE INTO roles (name, label) VALUES ('admin', 'Administrateur'), ('moderator', 'Modérateur'), ('ranger', 'Ranger'), ('user', 'Utilisateur'), ('blocked', 'Bloqué')`)
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
        INSERT OR IGNORE INTO user_roles (user_id, role_id)
        SELECT u.id, r.id
        FROM users u
        JOIN roles r ON r.name = 'admin'
        WHERE u.email = 'admin@dj.com'
    `)
    if err != nil {
        log.Println("Erreur seed admin role :", err)
    }
    _, err = db.Exec(`
        INSERT OR IGNORE INTO user_roles (user_id, role_id)
        SELECT u.id, r.id
        FROM users u
        JOIN roles r ON r.name = 'moderator'
        WHERE u.email = 'clever.girl@dj.com'
    `)
    if err != nil {
        log.Println("Erreur seed moderator role :", err)
    }
    _, err = db.Exec(`
        INSERT OR IGNORE INTO user_roles (user_id, role_id)
        SELECT u.id, r.id
        FROM users u
        JOIN roles r ON r.name = 'ranger'
        WHERE u.email = 'esther.stone@dj.com'
    `)
    if err != nil {
        log.Println("Erreur seed ranger role :", err)
    }
    _, err = db.Exec(`
        INSERT OR IGNORE INTO user_roles (user_id, role_id)
        SELECT u.id, r.id
        FROM users u
        JOIN roles r ON r.name = 'blocked'
        WHERE u.email = 'geek.osaure@dj.com'
    `)
    if err != nil {
        log.Println("Erreur seed blocked role :", err)
    }
    _, err = db.Exec(`
        INSERT OR IGNORE INTO user_roles (user_id, role_id)
        SELECT u.id, r.id
        FROM users u
        JOIN roles r ON r.name = 'user'
        WHERE u.email IN ('admin@dj.com', 'clever.girl@dj.com', 'esther.stone@dj.com', 'ben.pincus@dj.com', 'geek.osaure@dj.com')
    `)
    if err != nil {
        log.Println("Erreur seed user roles :", err)
    }

    _, err = db.Exec(`
        INSERT OR IGNORE INTO categories (name) VALUES ('Général'), ('Isla Nublar'), ('Herbivores'), ('Omnivores'), ('Carnivores'), ('Théories'), ('Incidents')`)
    if err != nil {
        log.Println("Erreur seed categories :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Dinosaures en liberté'
        FROM categories c, users u
        WHERE c.name = 'Général' AND u.username = 'admin'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Dinosaures en liberté')
    `)
    if err != nil {
        log.Println("Erreur seed topic général :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Nouvelles découvertes sur le parc'
        FROM categories c, users u
        WHERE c.name = 'Isla Nublar' AND u.username = 'PT1TEFUTEE86'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Nouvelles découvertes sur le parc')
    `)
    if err != nil {
        log.Println("Erreur seed topic Isla Nublar :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Les herbivores sont-ils sous-estimés ?'
        FROM categories c, users u
        WHERE c.name = 'Herbivores' AND u.username = 'Esth3rSt0nE'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Les herbivores sont-ils sous-estimés ?')
    `)
    if err != nil {
        log.Println("Erreur seed topic Herbivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Les plus dangereux omnivores'
        FROM categories c, users u
        WHERE c.name = 'Omnivores' AND u.username = 'BoBen'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Les plus dangereux omnivores')
    `)
    if err != nil {
        log.Println("Erreur seed topic Omnivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Qui est le vrai roi du parc ?'
        FROM categories c, users u
        WHERE c.name = 'Carnivores' AND u.username = 'G3ek05aure'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Qui est le vrai roi du parc ?')
    `)
    if err != nil {
        log.Println("Erreur seed topic Carnivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Théorie du retour des ADN perdus'
        FROM categories c, users u
        WHERE c.name = 'Théories' AND u.username = 'admin'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Théorie du retour des ADN perdus')
    `)
    if err != nil {
        log.Println("Erreur seed topic Théories :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Incident du tyrannosaure hier'
        FROM categories c, users u
        WHERE c.name = 'Incidents' AND u.username = 'PT1TEFUTEE86'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Incident du tyrannosaure hier')
    `)
    if err != nil {
        log.Println("Erreur seed topic Incidents :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Bienvenue sur ce topic dédié aux dinosaures en liberté. Partagez vos meilleures observations et vos questions.'
        FROM topics t, users u
        WHERE t.title = 'Dinosaures en liberté' AND u.username = 'admin'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Bienvenue sur ce topic dédié aux dinosaures en liberté.%')
    `)
    if err != nil {
        log.Println("Erreur seed post général :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'La nouvelle zone d"Isla Nublar semble être un terrain de jeu parfait pour les explorateurs. Qu"en pensez-vous ?'
        FROM topics t, users u
        WHERE t.title = 'Nouvelles découvertes sur le parc' AND u.username = 'PT1TEFUTEE86'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'La nouvelle zone d Isla Nublar semble être un terrain de jeu parfait%')
    `)
    if err != nil {
        log.Println("Erreur seed post Isla Nublar :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Je pense que les herbivores sont souvent sous-estimés. Ils ont une place cruciale dans l"équilibre du parc.'
        FROM topics t, users u
        WHERE t.title = 'Les herbivores sont-ils sous-estimés ?' AND u.username = 'Esth3rSt0nE'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Je pense que les herbivores sont souvent sous-estimés.%')
    `)
    if err != nil {
        log.Println("Erreur seed post Herbivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Les omnivores les plus dangereux sont ceux qui savent se fondre dans la nature. Quels sont vos favoris ?'
        FROM topics t, users u
        WHERE t.title = 'Les plus dangereux omnivores' AND u.username = 'BoBen'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Les omnivores les plus dangereux sont ceux qui savent se fondre%')
    `)
    if err != nil {
        log.Println("Erreur seed post Omnivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Ce sera un débat sans fin, mais je parie que le tyrannosaure règne en maître. Qui est selon vous le roi ?'
        FROM topics t, users u
        WHERE t.title = 'Qui est le vrai roi du parc ?' AND u.username = 'G3ek05aure'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Ce sera un débat sans fin, mais je parie que le tyrannosaure règne en maître.%')
    `)
    if err != nil {
        log.Println("Erreur seed post Carnivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'La théorie du retour des ADN perdus est fascinante : si on pouvait recréer des espèces éteintes, quels seraient les enjeux ?'
        FROM topics t, users u
        WHERE t.title = 'Théorie du retour des ADN perdus' AND u.username = 'admin'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'La théorie du retour des ADN perdus est fascinante%')
    `)
    if err != nil {
        log.Println("Erreur seed post Théories :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Un incident a été signalé avec le tyrannosaure hier soir. Quelqu"un a-t-il des infos supplémentaires ?'
        FROM topics t, users u
        WHERE t.title = 'Incident du tyrannosaure hier' AND u.username = 'PT1TEFUTEE86'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Un incident a été signalé avec le tyrannosaure hier soir.%')
    `)
    if err != nil {
        log.Println("Erreur seed post Incidents :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Photos et vidéos du parc'
        FROM categories c, users u
        WHERE c.name = 'Général' AND u.username = 'BoBen'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Photos et vidéos du parc')
    `)
    if err != nil {
        log.Println("Erreur seed topic photos :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Comportement social des herbivores'
        FROM categories c, users u
        WHERE c.name = 'Herbivores' AND u.username = 'PT1TEFUTEE86'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Comportement social des herbivores')
    `)
    if err != nil {
        log.Println("Erreur seed topic comportement herbivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Voici quelques photos récentes que j"ai prises lors de ma dernière patrouille.'
        FROM topics t, users u
        WHERE t.title = 'Photos et vidéos du parc' AND u.username = 'BoBen'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Voici quelques photos récentes%')
    `)
    if err != nil {
        log.Println("Erreur seed post photos :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Les herbivores forment souvent des groupes protecteurs autour des petits. J"ai observé ce comportement hier.'
        FROM topics t, users u
        WHERE t.title = 'Comportement social des herbivores' AND u.username = 'Esth3rSt0nE'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Les herbivores forment souvent des groupes protecteurs%')
    `)
    if err != nil {
        log.Println("Erreur seed post comportement herbivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content, parent_id)
        SELECT t.id, u.id, 'Super observation, merci pour les photos !', p.id
        FROM topics t, users u, posts p
        WHERE t.title = 'Photos et vidéos du parc' AND u.username = 'Esth3rSt0nE'
          AND p.topic_id = t.id AND p.user_id = (SELECT id FROM users WHERE username = 'BoBen')
          AND p.content LIKE 'Voici quelques photos récentes%'
          AND NOT EXISTS (SELECT 1 FROM posts WHERE parent_id = p.id AND content LIKE 'Super observation, merci pour les photos !')
    `)
    if err != nil {
        log.Println("Erreur seed reply photos :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content, parent_id)
        SELECT t.id, u.id, 'Est-ce que quelqu"un a remarqué des changements de migration cette saison ?', p.id
        FROM topics t, users u, posts p
        WHERE t.title = 'Comportement social des herbivores' AND u.username = 'G3ek05aure'
          AND p.topic_id = t.id AND p.user_id = (SELECT id FROM users WHERE username = 'Esth3rSt0nE')
          AND p.content LIKE 'Les herbivores forment souvent des groupes protecteurs%'
          AND NOT EXISTS (SELECT 1 FROM posts WHERE parent_id = p.id AND content LIKE 'Est-ce que quelqu un a remarqué des changements de migration cette saison ?')
    `)
    if err != nil {
        log.Println("Erreur seed reply herbivores :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Zones interdites et périmètres de sécurité'
        FROM categories c, users u
        WHERE c.name = 'Isla Nublar' AND u.username = 'admin'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Zones interdites et périmètres de sécurité')
    `)
    if err != nil {
        log.Println("Erreur seed topic zones interdites :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Des rangers ont récemment renforcé le périmètre près de la côte nord. Prudence recommandée.'
        FROM topics t, users u
        WHERE t.title = 'Zones interdites et périmètres de sécurité' AND u.username = 'admin'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Des rangers ont récemment renforcé le périmètre%')
    `)
    if err != nil {
        log.Println("Erreur seed post zones interdites :", err)
    }

    _, err = db.Exec(`
        INSERT INTO topics (category_id, user_id, title)
        SELECT c.id, u.id, 'Éthique de la recréation d"espèces'
        FROM categories c, users u
        WHERE c.name = 'Théories' AND u.username = 'BoBen'
        AND NOT EXISTS (SELECT 1 FROM topics WHERE title = 'Éthique de la recréation d espèces')
    `)
    if err != nil {
        log.Println("Erreur seed topic éthique :", err)
    }

    _, err = db.Exec(`
        INSERT INTO posts (topic_id, user_id, content)
        SELECT t.id, u.id, 'Si l"on recrée des espèces, devons-nous les protéger ou les étudier seulement ?'
        FROM topics t, users u
        WHERE t.title = 'Éthique de la recréation d espèces' AND u.username = 'admin'
        AND NOT EXISTS (SELECT 1 FROM posts WHERE topic_id = t.id AND user_id = u.id AND content LIKE 'Si l on recrée des espèces, devons-nous les protéger%')
    `)
    if err != nil {
        log.Println("Erreur seed post éthique :", err)
    }

    _, err = db.Exec(`
        INSERT INTO reports (reporter_id, target_type, target_id, reason, is_resolved)
        SELECT u.id, 'post', p.id, 'Contenu trop sensationnaliste', 0
        FROM users u, posts p
        WHERE u.username = 'PT1TEFUTEE86'
          AND p.content LIKE 'Bienvenue sur ce topic dédié aux dinosaures en liberté.%'
          AND NOT EXISTS (SELECT 1 FROM reports WHERE reporter_id = u.id AND target_type = 'post' AND target_id = p.id)
    `)
    if err != nil {
        log.Println("Erreur seed report 1 :", err)
    }

    _, err = db.Exec(`
        INSERT INTO reports (reporter_id, target_type, target_id, reason, is_resolved)
        SELECT u.id, 'post', p.id, 'Hors sujet et agressif', 0
        FROM users u, posts p
        WHERE u.username = 'Esth3rSt0nE'
          AND p.content LIKE 'Les omnivores les plus dangereux sont ceux qui savent se fondre dans la nature.%'
          AND NOT EXISTS (SELECT 1 FROM reports WHERE reporter_id = u.id AND target_type = 'post' AND target_id = p.id)
    `)
    if err != nil {
        log.Println("Erreur seed report 2 :", err)
    }

    _, err = db.Exec(`
        INSERT INTO reports (reporter_id, target_type, target_id, reason, is_resolved)
        SELECT u.id, 'post', p.id, 'Signalement urgent, possible fausse information', 0
        FROM users u, posts p
        WHERE u.username = 'BoBen'
          AND p.content LIKE 'Un incident a été signalé avec le tyrannosaure hier soir.%'
          AND NOT EXISTS (SELECT 1 FROM reports WHERE reporter_id = u.id AND target_type = 'post' AND target_id = p.id)
    `)
    if err != nil {
        log.Println("Erreur seed report 3 :", err)
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

    _, err = db.Exec(`
                INSERT INTO reports (reporter_id, target_type, target_id, reason, is_resolved)
                SELECT u.id, 'topic', t.id, 'Sujet encourage des expériences contraires à l''éthique', 0
        FROM users u, topics t
        WHERE u.username = 'Esth3rSt0nE' AND t.title = 'Éthique de la recréation d espèces'
          AND NOT EXISTS (SELECT 1 FROM reports WHERE reporter_id = u.id AND target_type = 'topic' AND target_id = t.id)
    `)
    if err != nil {
        log.Println("Erreur seed report topic :", err)
    }
}
