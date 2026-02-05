# Frontend - UpcycleConnect

Applications frontend développées en Vue.js.

## Objectif

Fournir des interfaces distinctes selon le type d’utilisateur :
- Particulier
- Professionnel / Artisan
- Salarié
- Administrateur (Backoffice)

Chaque application consomme la même API backend.

## Structure

- `apps/` : Applications Vue par rôle
  - `public-web/`
  - `pro-web/`
  - `staff-web/`
  - `backoffice/`
- `packages/` : Code partagé
  - `api-client/` : Client HTTP pour l’API
  - `ui/` : Composants UI communs
  - `auth/` : Gestion auth & permissions
  - `utils/` : Fonctions utilitaires

## Développement

Chaque app est une app Vue indépendante, mais partage du code via `packages/`.
