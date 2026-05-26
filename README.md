# Project Forum
### B1 - Evan Pertegaz, Louis Godard, Arthur Demarcq
 
## The Dark Jurassic
 
Dark Jurassic est un forum web développé en **Go** (backend) avec **HTML/CSS**. Il repose sur une base de données **SQLite** embarquée.
 
### Fonctionnalités principales
 
- Inscription, connexion et gestion de profil
- Catégories, topics et posts organisés hiérarchiquement
- Système de likes / dislikes
- Répondre aux posts
- Signalement de contenu et espace modération
---

### Prérequis
 
Il faut avoir **Go 1.22 ou supérieur** installé sur sa machine.
 
Pour vérifier :
 
```bash
go version
```

## Installation & Lancement

Pour commencer, on se rend dans le dossier contenant notre projet avec 

```
cd projet_forum/
```

Et on installe les dépendances puis on lance le serveur en une seule commande :
 
```bash
go run ./cmd/server/main.go
```
 
Le serveur démarre sur le port **8080** par défaut. On peut alors ouvrir le forum dans son navigateur à l'adresse depuis le terminal :
 
```
http://localhost:8080
```
