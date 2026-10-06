# ADR-003 : Gestion de l'authentification et des sessions

## Nomination courte

Utilisation de sessions pour gérer l'authentification des utilisateurs
et intégration de l'authentification via Google OAuth.

## Statut

Accepté

## Contexte

L'application est un forum dans lequel certaines fonctionnalités
nécessitent qu'un utilisateur soit authentifié.

Les utilisateurs doivent notamment pouvoir :

- créer un compte ;
- se connecter ;
- se déconnecter ;
- accéder à leur profil ;
- modifier leur profil ;
- supprimer leur compte ;
- effectuer certaines actions réservées aux utilisateurs authentifiés.

L'application doit donc pouvoir identifier un utilisateur entre
plusieurs requêtes HTTP.

Une solution d'authentification doit également permettre de conserver
l'état de connexion de l'utilisateur sans avoir à transmettre ses
informations d'identification à chaque requête.

Le projet propose également une authentification avec un compte Google.

## Options envisagées

### Option 1 : Utiliser uniquement les identifiants à chaque requête

L'utilisateur transmettrait ses identifiants lors de chaque requête
nécessitant une authentification.

**Avantages :**
- fonctionnement relativement simple à comprendre ;
- aucune gestion de session nécessaire.

**Inconvénients :**
- les identifiants devraient être transmis régulièrement ;
- moins pratique pour l'utilisateur ;
- mauvaise séparation entre l'authentification et les requêtes
  effectuées après la connexion ;
- augmente les risques liés à la transmission répétée des informations
  d'authentification.

### Option 2 : Utiliser des sessions côté serveur

Après une authentification réussie, une session est créée pour
l'utilisateur.

Les requêtes suivantes peuvent alors être associées à cette session
afin d'identifier l'utilisateur sans lui demander de se reconnecter
à chaque action.

**Avantages :**
- l'utilisateur reste connecté entre les différentes pages ;
- permet de gérer facilement la connexion et la déconnexion ;
- les données de session peuvent être gérées côté serveur ;
- adapté au fonctionnement d'une application web traditionnelle.

**Inconvénients :**
- nécessite de gérer la durée de vie des sessions ;
- nécessite de stocker et supprimer les sessions ;
- ajoute une gestion supplémentaire côté serveur.

### Option 3 : Utiliser des tokens JWT

L'utilisateur reçoit un token après son authentification et le transmet
ensuite lors des requêtes nécessitant une authentification.

**Avantages :**
- adapté aux architectures basées sur des API ;
- permet de transporter les informations nécessaires à
  l'authentification dans un token.

**Inconvénients :**
- gestion de l'expiration et de la révocation des tokens ;
- complexité supplémentaire pour le projet ;
- moins adapté aux besoins actuels d'une application utilisant
  principalement des pages HTML rendues côté serveur.

## Critères

La solution doit :

- permettre de maintenir la connexion d'un utilisateur ;
- identifier l'utilisateur lors des requêtes suivantes ;
- permettre une déconnexion ;
- être adaptée à une application web rendue côté serveur ;
- rester suffisamment simple à maintenir ;
- pouvoir être utilisée avec les fonctionnalités nécessitant
  une authentification.

## Décision finale

L'application utilise un système de sessions pour gérer
l'authentification des utilisateurs.

Une session est créée après une authentification réussie et permet
d'associer les requêtes suivantes à l'utilisateur connecté.

La gestion des sessions est séparée de la logique d'authentification
grâce aux composants dédiés du projet.

L'application propose également une authentification via Google OAuth.
Cette méthode permet à un utilisateur de se connecter à l'application
avec son compte Google.

Les routes dédiées à cette authentification sont notamment :

- `/auth/google/login`
- `/auth/google/callback`

## Conséquences

### Avantages

- l'utilisateur reste authentifié entre les différentes pages ;
- la déconnexion peut invalider la session ;
- les fonctionnalités nécessitant une authentification peuvent
  identifier l'utilisateur connecté ;
- la solution est adaptée au fonctionnement actuel du forum ;
- l'authentification Google permet de proposer une méthode de connexion
  supplémentaire.

### Inconvénients

- les sessions doivent être correctement gérées côté serveur ;
- les sessions doivent être supprimées ou invalidées lors de la
  déconnexion ;
- l'intégration de Google OAuth nécessite une configuration
  supplémentaire ;
- l'application dépend du fonctionnement du service Google pour
  l'authentification OAuth.