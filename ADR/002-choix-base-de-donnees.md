# ADR-002 : Choix de SQLite comme système de base de données

## Nomination courte

Utilisation de SQLite pour le stockage des données de l'application.

## Statut

Accepté

## Contexte

L'application est un forum développé en Go. Elle doit stocker
différentes données nécessaires au fonctionnement de l'application,
notamment :

- les utilisateurs ;
- les rôles ;
- les catégories ;
- les sujets ;
- les publications ;
- les likes ;
- les signalements ;
- les sessions ;
- les notifications.

Le projet nécessite donc une base de données permettant de stocker
ces informations de manière persistante.

L'application étant principalement destinée à fonctionner comme un
projet de taille limitée, la solution doit rester simple à installer,
à configurer et à utiliser.

## Options envisagées

### Option 1 : SQLite

SQLite est une base de données légère fonctionnant directement avec
l'application, sans nécessiter de serveur de base de données séparé.

**Avantages :**
- installation et configuration simples ;
- aucune infrastructure supplémentaire nécessaire ;
- base de données stockée dans un fichier ;
- adaptée à un projet de taille limitée ;
- intégration simple avec une application Go.

**Inconvénients :**
- moins adaptée à une application nécessitant une très forte
  charge ou un grand nombre de connexions simultanées ;
- les possibilités de mise à l'échelle sont plus limitées qu'avec
  une base de données serveur.

### Option 2 : PostgreSQL

PostgreSQL est une base de données relationnelle complète pouvant
supporter des applications importantes et une charge élevée.

**Avantages :**
- très bonnes performances ;
- gestion avancée des transactions ;
- adaptée aux applications évolutives ;
- bonne gestion des connexions simultanées.

**Inconvénients :**
- nécessite un serveur de base de données ;
- configuration et déploiement plus complexes ;
- infrastructure plus importante pour les besoins actuels du projet.

### Option 3 : MySQL

MySQL est également une base de données relationnelle largement
utilisée.

**Avantages :**
- bonne maturité ;
- bonnes performances ;
- adaptée aux applications web.

**Inconvénients :**
- nécessite également une infrastructure de base de données ;
- configuration plus complexe que SQLite ;
- cette complexité n'est pas nécessaire pour les besoins actuels
  du projet.

## Critères

Le choix de la base de données doit prendre en compte :

- la simplicité d'installation ;
- la facilité de configuration ;
- l'intégration avec l'application Go ;
- les besoins actuels du projet ;
- les performances nécessaires ;
- la possibilité de faire évoluer l'application.

## Décision finale

SQLite est retenue comme système de gestion de base de données
pour l'application.

Cette solution permet de conserver une architecture simple tout en
disposant d'une base de données relationnelle adaptée aux besoins
actuels du forum.

La base de données est gérée par l'application et les changements
de structure sont pris en charge par le système de migrations présent
dans `internal/database`.

## Conséquences

### Avantages

- configuration simplifiée de l'application ;
- aucune base de données externe à installer pour exécuter le projet ;
- déploiement facilité ;
- stockage persistant des données ;
- solution adaptée à la taille actuelle du projet.

### Inconvénients

- limitations pour une montée en charge importante ;
- moins adaptée à un environnement avec de nombreuses connexions
  simultanées ;
- un changement de base de données pourrait être nécessaire si
  l'application devait évoluer vers une architecture à plus grande
  échelle.