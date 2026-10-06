# ADR-006 : Gestion des migrations de la base de données

## Nomination courte

Utilisation d'un système de migrations pour gérer l'évolution du
schéma de la base de données.

## Statut

Accepté

## Contexte

L'application utilise une base de données SQLite pour stocker les
informations nécessaires au fonctionnement du forum.

Au cours du développement, la structure de la base de données peut
évoluer. De nouvelles tables, colonnes ou contraintes peuvent être
nécessaires lorsque de nouvelles fonctionnalités sont ajoutées.

Il est donc nécessaire de disposer d'un mécanisme permettant de
modifier le schéma de la base de données de manière contrôlée et
reproductible.

Le projet possède un système de migrations situé dans le dossier
`internal/database/`.

## Options envisagées

### Option 1 : Modifier directement la base de données

Les modifications du schéma seraient effectuées manuellement
directement dans la base de données.

**Avantages :**
- méthode simple pour une modification ponctuelle ;
- aucune logique de migration supplémentaire à développer.

**Inconvénients :**
- historique des modifications difficile à conserver ;
- risque d'oublier une modification lors de l'installation du projet ;
- difficile de reproduire exactement la même structure sur plusieurs
  environnements ;
- risque d'erreurs lors des modifications successives.

### Option 2 : Recréer la base de données à chaque modification

La base de données pourrait être supprimée puis recréée avec
la nouvelle structure.

**Avantages :**
- structure de la base facilement reproductible ;
- implémentation relativement simple.

**Inconvénients :**
- perte des données existantes ;
- impossible à utiliser correctement lorsque des données doivent
  être conservées ;
- peu adapté à l'évolution progressive d'une application.

### Option 3 : Utiliser un système de migrations

Chaque modification du schéma est enregistrée dans une migration.

Les migrations sont exécutées afin de faire évoluer progressivement
la structure de la base de données.

**Avantages :**
- conservation de l'historique des modifications ;
- évolution progressive du schéma ;
- conservation des données existantes ;
- possibilité de reproduire la structure de la base ;
- meilleure organisation des changements liés à la base de données.

**Inconvénients :**
- nécessite de maintenir les migrations ;
- une migration incorrecte peut provoquer une erreur lors de
  l'initialisation de l'application ;
- ajoute une couche de gestion supplémentaire.

## Critères

La solution doit :

- permettre de faire évoluer la structure de la base de données ;
- conserver les données existantes ;
- permettre de reproduire la structure de la base ;
- faciliter le suivi des modifications ;
- être adaptée à l'utilisation de SQLite ;
- rester suffisamment simple pour le projet.

## Décision finale

Un système de migrations est utilisé pour gérer l'évolution du schéma
de la base de données.

Les migrations sont regroupées dans le dossier
`internal/database/` et sont exécutées lors de l'initialisation
de l'application.

Cette approche permet d'appliquer progressivement les modifications
nécessaires à la base de données sans devoir supprimer et recréer
celle-ci à chaque évolution du projet.

## Conséquences

### Avantages

- évolution contrôlée de la base de données ;
- conservation des données existantes ;
- historique des modifications du schéma ;
- structure de la base reproductible ;
- meilleure organisation des évolutions de la base de données.

### Inconvénients

- nécessité de maintenir les migrations avec le reste du projet ;
- les migrations doivent être exécutées dans le bon ordre ;
- une erreur dans une migration peut empêcher l'initialisation
  correcte de la base de données ;
- les anciennes migrations doivent être conservées afin de pouvoir
  reconstruire correctement la structure de la base.