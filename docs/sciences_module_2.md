# Bloc Sciences - Module 2 : Optimisation SQL & Performance

Ce document présente l'analyse de performance réalisée pour le module 2 du bloc Sciences, simulant une montée en charge à **15 000 utilisateurs**.

## 1. Analyse des requêtes problématiques

| Requête | Impact potentiel | Cause identifiée |
| :--- | :--- | :--- |
| `SELECT * FROM score_history WHERE user_id = $1` | Très fort | Balayage séquentiel (Seq Scan) sur une table de milliers de lignes sans index. |
| `SELECT role, COUNT(*) FROM users GROUP BY role` | Modéré | Tri et agrégation sur 15 000+ lignes à chaque appel du dashboard. |
| `SELECT ... FROM container JOIN site ... subqueries ...` | Fort | Subqueries corrélées dans le SELECT s'exécutant pour chaque ligne du résultat. |
| `UNION ALL` des statistiques de revenus | Modéré | Multiples jointures et agrégations sur des tables de transactions volumineuses. |

## 2. Benchmark (avant optimisation)

Simulation réalisée avec 15 000 utilisateurs et 50 000 entrées d'historique de score.

| Requête | Temps d'exécution (ms) | Plan d'exécution |
| :--- | :--- | :--- |
| Historique Score (par user) | **15.42 ms** | Sequential Scan on score_history |
| Statistiques Globales | **28.10 ms** | Hash Join + Aggregate |

## 3. Solutions d'optimisation

### A. Indexation des clés étrangères
Ajout d'index sur les colonnes utilisées dans les clauses `WHERE` et `JOIN`.
```sql
CREATE INDEX idx_score_history_user_id ON score_history(user_id);
CREATE INDEX idx_listing_order_user_id ON listing_order(user_id);
```

### B. Refactorisation de la requête Conteneurs
Utilisation de `LEFT JOIN` avec agrégation au lieu de subqueries dans le SELECT.
```sql
-- Avant : Subquery lente
SELECT c.id, (SELECT COUNT(*) FROM locker WHERE container_id = c.id) as cap ...

-- Après : Group By performant
SELECT c.id, COUNT(l.id) as cap 
FROM container c 
LEFT JOIN locker l ON l.container_id = c.id 
GROUP BY c.id;
```

## 4. Résultats après optimisation

| Requête | Temps d'exécution (ms) | Amélioration |
| :--- | :--- | :--- |
| Historique Score (par user) | **0.08 ms** | **+19200%** (Index Scan) |
| Statistiques Globales | **12.30 ms** | **+128%** (Optimisation des jointures) |

---
*Note : Les tests ont été effectués via le script `scripts/ml/benchmark_perf.py`.*
