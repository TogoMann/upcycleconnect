# Router

Ce dossier gère la configuration du routeur Gin.

## Rôle

- Initialiser l’instance Gin
- Déclarer les routes globales
- Enregistrer les routes de chaque module
- Appliquer les middlewares (auth, rôles, CORS, etc.)

## Exemples de responsabilités

- Séparer les routes publiques et protégées
- Grouper les routes par rôle (admin, staff, pro, user)
- Brancher les routes des modules métier

## Objectif

Centraliser la définition des routes et garantir une structure claire de l’API.
