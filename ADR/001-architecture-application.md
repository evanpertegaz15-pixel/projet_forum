# ADR-001 : Architecture en couches de l'application

## Nomination courte

Adoption d'une architecture en couches pour l'application du forum.

## Statut

Accepté

## Contexte

L'application est un forum développé en Go.
Elle doit gérer plusieurs fonctionnalités telles que :

- l'authentification des utilisateurs ;
- les profils ;
- les catégories et sujets ;
- les publications et réponses ;
- les likes ;
- les signalements ;
- la gestion des utilisateurs ;
- les notifications.

Le projet doit rester maintenable et permettre de faire évoluer
les différentes fonctionnalités sans mélanger la logique HTTP,
la logique métier et l'accès aux données.

## Options envisagées

### Option 1 : Tout regrouper dans les handlers

**Avantages :**
- développement rapide ;
- structure simple pour un petit projet.

**Inconvénients :**
- logique métier mélangée avec HTTP ;
- fichiers difficiles à maintenir ;
- code plus difficile à tester ;
- risque de duplication.

### Option 2 : Utiliser une architecture en couches

Séparer l'application en plusieurs responsabilités :

- `handlers` : gestion des requêtes HTTP ;
- `services` : logique métier ;
- `models` : accès et représentation des données ;
- `middleware` : traitements appliqués aux requêtes ;
- `utils` : fonctionnalités communes ;
- `templates` : présentation.

**Avantages :**
- responsabilités clairement séparées ;
- meilleure maintenabilité ;
- code plus facilement réutilisable ;
- possibilité de faire évoluer une couche sans modifier les autres.

**Inconvénients :**
- davantage de fichiers ;
- architecture plus complexe pour une petite fonctionnalité.

## Critères

Le choix doit permettre :

- de séparer les responsabilités ;
- de faciliter la maintenance ;
- de limiter les dépendances entre les différentes parties ;
- de faciliter l'évolution du projet.

## Décision finale

L'application utilise une architecture en couches.

Les handlers sont responsables des requêtes HTTP et délèguent
la logique métier aux services. Les services utilisent les modèles
pour accéder aux données.

Cette organisation permet notamment de conserver les handlers
relativement indépendants de la logique d'accès à la base de données.

## Conséquences

### Avantages

- meilleure organisation du projet ;
- séparation claire des responsabilités ;
- maintenance facilitée ;
- ajout de nouvelles fonctionnalités facilité.

### Inconvénients

- davantage de fichiers et de code de liaison ;
- certaines fonctionnalités simples nécessitent plusieurs couches ;
- la structure demande de respecter les responsabilités de chaque couche.