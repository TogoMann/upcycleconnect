# Auth

Ce package gère toute la logique d’authentification et d’autorisation côté frontend.

## Rôle

- Stockage et gestion du token (JWT)
- Récupération de l’utilisateur courant
- Vérification des rôles et permissions
- Guards de routes (protéger certaines pages)
- Helpers du type `isAdmin()`, `isStaff()`, etc.

## Objectif

Centraliser la sécurité côté frontend et garantir un comportement cohérent
entre toutes les applications (public, pro, staff, admin).
