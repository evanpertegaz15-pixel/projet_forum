# ADR-005 : Mise en place d'un système de Rate Limiting

## Nomination courte

Limitation du nombre de requêtes par adresse IP afin de protéger
l'application contre les requêtes excessives.

## Statut

Accepté

## Contexte

L'application est un forum accessible via un serveur HTTP.

Un utilisateur ou un client peut envoyer un nombre important de
requêtes au serveur sur une courte période. Un nombre excessif de
requêtes peut entraîner une utilisation inutile des ressources du
serveur et dégrader les performances de l'application.

Il est donc nécessaire de mettre en place un mécanisme permettant
de limiter le nombre de requêtes effectuées sur une période donnée.

Le système doit également fonctionner pour les utilisateurs
authentifiés et non authentifiés.

## Options envisagées

### Option 1 : Ne pas mettre en place de limitation

Toutes les requêtes reçues par le serveur seraient traitées sans
limitation.

**Avantages :**
- aucune logique supplémentaire à implémenter ;
- fonctionnement simple.

**Inconvénients :**
- aucune protection contre un nombre excessif de requêtes ;
- risque de dégradation des performances ;
- possibilité d'abus du serveur.

### Option 2 : Limitation par adresse IP

Chaque adresse IP possède son propre limiteur de requêtes.

Lorsqu'une adresse dépasse la limite autorisée, les requêtes
supplémentaires sont refusées temporairement.

**Avantages :**
- solution simple à mettre en place ;
- fonctionne également pour les utilisateurs non authentifiés ;
- ne nécessite pas de modification de la base de données ;
- permet de limiter les requêtes provenant d'un même client.

**Inconvénients :**
- plusieurs utilisateurs partageant une même adresse IP peuvent
  partager la même limite ;
- l'adresse IP ne permet pas toujours d'identifier précisément
  un utilisateur ;
- les utilisateurs utilisant certains réseaux ou proxys peuvent
  être concernés par une même limitation.

### Option 3 : Limitation par utilisateur authentifié

La limite serait associée directement au compte utilisateur.

**Avantages :**
- limitation plus précise pour les utilisateurs authentifiés ;
- permet de définir des limites différentes selon les utilisateurs.

**Inconvénients :**
- ne protège pas directement les utilisateurs non authentifiés ;
- nécessite une gestion supplémentaire ;
- nécessite de prendre en compte les utilisateurs anonymes.

## Critères

La solution doit :

- limiter les requêtes excessives ;
- fonctionner pour les utilisateurs authentifiés et non authentifiés ;
- être simple à intégrer au serveur existant ;
- avoir un faible impact sur les performances ;
- ne pas nécessiter de modification supplémentaire de la base de données.

## Décision finale

Un système de Rate Limiting basé sur l'adresse IP est utilisé.

Chaque adresse IP possède son propre limiteur de requêtes.

Le mécanisme est implémenté à l'aide du package
`golang.org/x/time/rate`.

La configuration actuelle autorise un burst initial de 10 requêtes
et permet ensuite environ une requête par seconde.

Lorsqu'une adresse IP dépasse cette limite, la requête est refusée
avec le code HTTP `429 Too Many Requests`.

Une page dédiée au code `429` est affichée afin d'informer
explicitement l'utilisateur qu'il a effectué trop de requêtes et
qu'il doit patienter avant de réessayer.

Les informations concernant les visiteurs inactifs sont également
nettoyées périodiquement afin d'éviter de conserver indéfiniment
leurs limiteurs en mémoire.

## Conséquences

### Avantages

- protection contre les requêtes excessives ;
- fonctionnement pour les utilisateurs connectés et anonymes ;
- aucune modification nécessaire de la base de données ;
- implémentation relativement simple ;
- réponse HTTP `429` clairement identifiable ;
- page d'erreur dédiée permettant d'informer l'utilisateur.

### Inconvénients

- plusieurs utilisateurs partageant une même adresse IP partagent
  également la même limite ;
- une adresse IP ne représente pas nécessairement un utilisateur
  unique ;
- les limites actuelles peuvent nécessiter un ajustement en fonction
  de l'évolution du trafic ;
- les informations des visiteurs sont conservées temporairement
  en mémoire.