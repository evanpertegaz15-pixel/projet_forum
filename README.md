# Projet Forum
### B1 - Evan Pertegaz, Louis Godard, Arthur Demarcq
 
## The Dark Jurassic

### Présentation

Le Dark Jurassic s'inspire du [site](https://jurassicpark.fandom.com/fr/wiki/Dark_Jurassic) du même nom dans la série télévisée Netflix [Jurassic World : La Théorie du Chaos](https://jurassicpark.fandom.com/fr/wiki/Jurassic_World_:_La_Th%C3%A9orie_du_Chaos).

Dark Jurassic est un forum web développé en **Go** (backend) avec **HTML/CSS**. Il repose sur une base de données **SQLite** embarquée.
 
### Fonctionnalités principales
 
- Inscription, connexion et gestion de profil
- Catégories, topics et posts organisés hiérarchiquement
- Système de likes / dislikes
- Répondre aux posts
- Signalement de contenu et espace modération
---

### Prérequis
 
Il faut avoir **Go 1.25** installé sur sa machine.
 
Pour vérifier :
 
```bash
go version
```

## Installation & Lancement

### Installation
- Cloner le dépôt
  - Dans un terminal, utiliser `git clone <lien_repo>` depuis le dossier où l'on souhaite l'enregistrer.

### Lancement
 - Pour commencer, on se rend dans le dossier contenant notre projet avec :
    ```bash
    cd projet_forum/
    ```

- Et on installe les dépendances puis on lance le serveur en une seule commande :
    ```bash
    go run ./cmd/server/main.go
    ```
 
- Le serveur démarre sur le port **8080** par défaut. On peut alors ouvrir le forum dans son navigateur à l'adresse depuis le terminal :
    ```pwsh
    http://localhost:8080
    ```
