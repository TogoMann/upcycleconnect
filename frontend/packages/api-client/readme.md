# API Client

Ce package contient le client HTTP utilisé pour communiquer avec l’API backend.

## Rôle

- Centraliser les appels HTTP (fetch / axios)
- Gérer les headers (Authorization, Content-Type, etc.)
- Gérer les erreurs globales (401, 403, 500, etc.)
- Fournir des fonctions prêtes à l’emploi pour les modules frontend

## Exemples

- `getCurrentUser()`
- `login()`
- `createAnnonce()`
- `listOffres()`

## Objectif

Avoir un point d’entrée unique pour l’API et éviter la duplication de logique réseau
dans chaque application frontend.
