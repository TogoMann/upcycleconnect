# Backend - API UpcycleConnect

API REST développée en Go avec le framework Gin.

## Responsabilités

- Authentification et autorisation (JWT, rôles)
- Gestion des utilisateurs (particuliers, pros, salariés, admins)
- Gestion des annonces, offres, événements, conteneurs, etc.
- Accès base de données
- Sécurité et validation des requêtes

## Structure

- `cmd/api/` : Point d’entrée de l’application
- `internal/` : Code métier de l’API
  - `config/` : Configuration
  - `database/` : Connexion et gestion DB
  - `middlewares/` : Middlewares Gin (auth, rôles, CORS, logs, etc.)
  - `modules/` : Modules métier (users, annonces, offres, etc.)
  - `router/` : Déclaration des routes
  - `utils/` : Fonctions utilitaires
- `migrations/` : Scripts de migration DB
- `tests/` : Tests automatisés

## Lancer en local (hors Docker)

```bash
go run ./cmd/api
```

Écoute sur le port 8080
