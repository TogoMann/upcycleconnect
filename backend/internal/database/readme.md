# Database

Ce dossier gère tout ce qui concerne l’accès à la base de données.

## Rôle

- Initialiser la connexion à la base (PostgreSQL / MariaDB)
- Gérer le pool de connexions
- Fournir l’accès DB aux repositories
- Gérer éventuellement les transactions

## Exemples de contenu

- `db.go` : Initialisation de la connexion
- Helpers pour vérifier l’état de la DB
- Gestion des erreurs de connexion

## Objectif

Centraliser l’accès à la base de données et éviter les connexions
dupliquées ou mal configurées dans les modules métier.
