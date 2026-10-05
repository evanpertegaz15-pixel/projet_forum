# Projet Forum
### B1 - Evan Pertegaz, Louis Godard, Arthur Demarcq
 
## The Dark Jurassic

### Présentation

Le Dark Jurassic s'inspire du [site](https://jurassicpark.fandom.com/fr/wiki/Dark_Jurassic) du même nom dans la série télévisée Netflix [Jurassic World : La Théorie du Chaos](https://jurassicpark.fandom.com/fr/wiki/Jurassic_World_:_La_Th%C3%A9orie_du_Chaos). Il est destiné aux fans de la saga, et leur permet de se retrouver autour de sujets communs.

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

## Démonstration

- En arrivant sur le site, vous pouvez consulter les derniers posts dans le fil d'actualité. Ainsi qu'accéder aux catégories tendances.
- Si vous désirez pouvoir interagir avec les autres utilisateurs, il faut créer un compte.
  - Pour ce faire, cliquer en haut à droite de l'écran sur le logo de connexion / inscription et suivre les indications.
- Une fois connecté, vous pouvez vous rendre dans la partie catégorie et topics et commencer à créer des posts.
- Si toutefois vous voulez personnaliser votre profil, vous pouvez vous rendre sur votre page de la même manière que pour se connecter.
  - De là, appuyer sur "Modifier le profil", et n'oubliez pas d'enregistrer !

## Structure du dépôt

### Racine

La racine contient :
- Les fichiers de configuration `go.mod` et `go.sum`
- La base de données `forum.db`, quand le projet est lancé

---

### Backend

Le backend se trouve dans `/internal` :
- Le fichier de configuration `config.go` est dans `/config` et permet de relier les différents services entre eux
- Pour la base de données : `/database` et `/models`
- La gestion des routes, et interactions est principalement dans : `handlers`, `services` et `utils`

---

### Frontend

Le frontend se trouve dans :
- `/internal/templates` pour les pages **HTML** en elles-mêmes
- `/static` pour les assets tels que les images, ou le **CSS** pour le style