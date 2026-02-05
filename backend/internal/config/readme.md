# Config

Ce dossier contient la gestion de la configuration de l’application backend.

## Rôle

- Charger les variables d’environnement
- Centraliser les paramètres (ports, DB, JWT, mode debug, etc.)
- Fournir une configuration typée au reste de l’application

## Exemples de contenu

- `config.go` : Structure de configuration globale
- Chargement depuis `.env` ou variables système
- Validation des valeurs critiques (ex: secrets, URLs, ports)

## Objectif

Éviter les valeurs “en dur” dans le code et faciliter le déploiement
en environnement de développement, test et production.
