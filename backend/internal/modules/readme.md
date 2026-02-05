# Modules métier

Chaque dossier ici représente un module fonctionnel de l’application.

Exemples :
- `users/` : Gestion des utilisateurs
- `auth/` : Authentification
- `annonces/` : Annonces de dons / ventes
- `offres/` : Offres, services, formations
- `conteneurs/` : Gestion des dépôts d’objets
- `evenements/` : Événements et ateliers

## Structure type d’un module

- `handler.go` : Handlers HTTP (controllers)
- `service.go` : Logique métier
- `repository.go` : Accès base de données
- `model.go` : Structures de données
- `routes.go` : Déclaration des routes du module

Ce découpage permet de garder un code clair, testable et maintenable.
