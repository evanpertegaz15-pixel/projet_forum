# ADR-004 : Rendu des pages HTML côté serveur

## Nomination courte

Utilisation des templates HTML Go pour générer les pages de l'application
côté serveur.

## Statut

Accepté

## Contexte

L'application est un forum développé en Go qui doit afficher plusieurs
pages accessibles depuis un navigateur, notamment :

- la page d'accueil ;
- les catégories ;
- les sujets ;
- les publications ;
- les profils ;
- les pages de connexion et d'inscription ;
- les pages d'erreur.

Le projet doit choisir une méthode permettant de générer et d'afficher
ces pages tout en conservant une architecture adaptée à sa taille.

L'application utilise actuellement des templates HTML stockés dans
`internal/templates/` et rendus par le serveur Go.

## Options envisagées

### Option 1 : Rendu HTML côté serveur avec les templates Go

Les pages HTML sont générées directement par le serveur Go à partir
de templates.

Les données nécessaires à l'affichage sont transmises par les handlers
aux templates.

**Avantages :**
- architecture simple ;
- intégration directe avec le serveur Go ;
- pas besoin d'un frontend séparé ;
- peu de dépendances supplémentaires ;
- les pages peuvent être générées directement à partir des données
  récupérées par l'application.

**Inconvénients :**
- une partie de l'affichage dépend du serveur ;
- les interactions dynamiques nécessitent du JavaScript supplémentaire ;
- les templates peuvent devenir plus complexes lorsque les interfaces
  deviennent importantes.

### Option 2 : Utiliser un framework frontend

L'interface pourrait être développée avec un framework frontend comme
React, Vue ou Angular.

Le serveur Go fournirait alors principalement les données nécessaires
au frontend.

**Avantages :**
- permet de créer des interfaces très dynamiques ;
- séparation plus importante entre le frontend et le backend ;
- possibilité de réutiliser des composants frontend.

**Inconvénients :**
- architecture plus complexe ;
- nécessité de maintenir une application frontend séparée ;
- davantage de dépendances ;
- temps de développement plus important ;
- solution disproportionnée par rapport aux besoins actuels du projet.

### Option 3 : Utiliser une API avec un frontend séparé

Le serveur Go pourrait uniquement fournir une API et une application
frontend indépendante serait responsable de l'affichage.

**Avantages :**
- séparation claire entre backend et frontend ;
- API réutilisable par différentes applications ;
- possibilité de développer plusieurs clients.

**Inconvénients :**
- complexité supplémentaire ;
- nécessité de maintenir deux parties distinctes ;
- gestion supplémentaire de l'authentification et des échanges
  entre le frontend et le backend ;
- solution plus importante que nécessaire pour le projet actuel.

## Critères

La solution doit :

- être simple à mettre en place ;
- s'intégrer facilement avec le serveur Go ;
- limiter les dépendances ;
- permettre de créer les différentes pages du forum ;
- faciliter la maintenance du projet ;
- être adaptée à la taille et aux besoins actuels de l'application.

## Décision finale

Le rendu des pages HTML est réalisé côté serveur à l'aide des templates
Go.

Les templates sont regroupés dans `internal/templates/` et les handlers
transmettent les données nécessaires à leur génération.

Le rendu est centralisé à travers une fonction dédiée permettant de
charger et d'exécuter les templates.

Cette approche permet de conserver une architecture simple avec un
serveur Go responsable à la fois du traitement des requêtes et de la
génération des pages HTML.

## Conséquences

### Avantages

- architecture simple et cohérente avec le projet ;
- aucun frontend séparé à maintenir ;
- nombre de dépendances limité ;
- intégration directe avec les handlers et les services Go ;
- développement et déploiement simplifiés ;
- les templates peuvent être réutilisés pour les différentes pages
  de l'application.

### Inconvénients

- les interfaces très dynamiques nécessitent du JavaScript
  supplémentaire ;
- le rendu des pages dépend du serveur ;
- une évolution vers une application frontend complète nécessiterait
  une modification plus importante de l'architecture.