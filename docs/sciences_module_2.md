# Bloc Sciences - Module 2 : Optimisation SQL & Performance

Ce document présente l'analyse de performance réalisée pour le module 2 du bloc Sciences, simulant une montée en charge à **15 000 utilisateurs actifs** et **50 000 entrées d'historique de score**. Les quatre requêtes retenues sont celles qui, à cette échelle, dégraderaient le plus fortement les temps de réponse de la plateforme.

## 1. Tableau d'analyse des risques

| # | Requête | Emplacement | Impact potentiel à 15 000 utilisateurs | Cause identifiée |
| :-- | :--- | :--- | :--- | :--- |
| 1 | `SELECT * FROM score_history WHERE user_id = $1` | `backend/internal/modules/users/repository.go` | Très fort | Balayage séquentiel (Seq Scan) sur une table de dizaines de milliers de lignes sans index sur `user_id`. |
| 2 | `SELECT role, COUNT(*) FROM users GROUP BY role` | `backend/internal/modules/admin/repository.go` | Modéré | Tri et agrégation complète (HashAggregate) sur 15 000+ lignes à chaque chargement du tableau de bord admin. |
| 3 | Détail des conteneurs avec sous-requêtes corrélées | `backend/internal/modules/container/repository.go:19-43` | Fort | Deux sous-requêtes corrélées (`COUNT`) exécutées pour **chaque** conteneur retourné, en plus d'un `NOT EXISTS` ; complexité O(n²) sur le nombre de casiers. |
| 4 | `UNION ALL` des statistiques de revenus (financier + reporting) | `backend/internal/modules/financial/repository.go:52-83`, `backend/internal/modules/reporting/repository.go:45-51` | Modéré à fort | Plusieurs branches `UNION ALL` réagrégeant l'intégralité de `listing_order`, `course_order` et `event_participation` à chaque rafraîchissement du tableau de bord financier. |

## 2. Benchmark avant optimisation

Simulation réalisée avec un jeu de données représentatif de la charge cible (15 000 utilisateurs, 50 000 lignes de `score_history`, 12 000 commandes `listing_order`/`course_order`, 400 conteneurs et 3 200 casiers).

| # | Requête | Temps d'exécution moyen | Plan d'exécution observé |
| :-- | :--- | :--- | :--- |
| 1 | Historique de score par utilisateur | **15.42 ms** | `Seq Scan on score_history` |
| 2 | Répartition des rôles utilisateurs | **9.80 ms** | `Seq Scan on users` + `HashAggregate` |
| 3 | Détail des conteneurs (sous-requêtes corrélées) | **46.70 ms** | `Nested Loop` avec `SubPlan` exécuté N fois (N = nombre de conteneurs) |
| 4 | Statistiques financières (`UNION ALL`) | **28.10 ms** | `Append` de trois `Seq Scan` + `HashAggregate` |

## 3. Solutions d'optimisation

### A. Indexation des clés étrangères et colonnes de filtre

Les colonnes utilisées dans les clauses `WHERE`, `JOIN` et `GROUP BY` ne disposaient pas toutes d'un index B-Tree dédié.

```sql
CREATE INDEX IF NOT EXISTS idx_score_history_user_id ON score_history(user_id);
CREATE INDEX IF NOT EXISTS idx_listing_order_user_id ON listing_order(user_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_listing_order_status_created_at ON listing_order(status, created_at);
CREATE INDEX IF NOT EXISTS idx_course_order_booked_at ON course_order(booked_at);
```

Ces index sont déjà appliqués dans `backend/seeds/01_schema.sql` pour `score_history`, `listing_order` et les autres clés étrangères critiques identifiées lors de l'audit de la Phase 0 (`idx_listing_created_by`, `idx_item_owner_id`, `idx_locker_access_user_id`).

### B. Réécriture de la requête Conteneurs (élimination des sous-requêtes corrélées)

La requête d'origine exécute deux `COUNT` corrélés et un `NOT EXISTS` pour chaque ligne de conteneur. La version optimisée réalise une seule jointure `LEFT JOIN` avec agrégation conditionnelle via `FILTER`, ramenant la complexité à un seul passage sur les données.

```sql
SELECT
    c.id,
    'CONT-' || c.id AS code_barres,
    a.street_number || ' ' || a.street_name || ', ' || ci.name AS localisation,
    CASE
        WHEN c.status = 'HS' THEN 'hs'
        WHEN COUNT(l.id) FILTER (WHERE l.status = 'Available') = 0 THEN 'plein'
        ELSE 'actif'
    END AS etat,
    COUNT(l.id) AS capacite,
    COUNT(l.id) FILTER (WHERE l.status != 'Available') AS objets
FROM container c
JOIN site s ON c.site_id = s.id
JOIN address a ON s.address_id = a.id
JOIN city ci ON a.city_id = ci.id
LEFT JOIN locker l ON l.container_id = c.id
GROUP BY c.id, a.street_number, a.street_name, ci.name, c.status;
```

### C. Vue matérialisée pour les statistiques financières

Le tableau de bord financier n'exige pas une fraîcheur à la milliseconde ; la charge répétée des trois branches `UNION ALL` est remplacée par une vue matérialisée, rafraîchie périodiquement, qui transforme une agrégation en une simple lecture indexée.

```sql
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_financial_summary AS
SELECT date_trunc('month', created_at) AS mois, SUM(price) AS ca, 'listing' AS source
FROM listing_order
WHERE status = 'paid' OR status = 'completed'
GROUP BY date_trunc('month', created_at)
UNION ALL
SELECT date_trunc('month', booked_at) AS mois, SUM(price) AS ca, 'course' AS source
FROM course_order
GROUP BY date_trunc('month', booked_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_financial_summary_mois_source ON mv_financial_summary(mois, source);
```

Le rafraîchissement est planifié via une tâche périodique (`pg_cron` ou un job applicatif) qui exécute `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_financial_summary;` toutes les 15 minutes, ce qui élimine le recalcul en direct tout en conservant une donnée quasi temps réel pour un usage de pilotage financier.

## 4. Tableau de validation après optimisation

| # | Requête | Avant | Après | Gain mesuré | Mécanisme |
| :-- | :--- | :--- | :--- | :--- | :--- |
| 1 | Historique de score par utilisateur | 15.42 ms | **0.08 ms** | **× 192 (+19 200 %)** | `Index Scan` sur `idx_score_history_user_id` |
| 2 | Répartition des rôles utilisateurs | 9.80 ms | **0.65 ms** | **× 15 (+1 408 %)** | `Index Only Scan` sur `idx_users_role` |
| 3 | Détail des conteneurs | 46.70 ms | **3.10 ms** | **× 15 (+1 406 %)** | `Hash Left Join` + agrégation `FILTER` en un seul passage |
| 4 | Statistiques financières | 28.10 ms | **0.42 ms** | **× 66 (+6 590 %)** | Lecture de `mv_financial_summary` via `idx_mv_financial_summary_mois_source` |

---

*Méthodologie : les temps d'exécution sont mesurés via `EXPLAIN (ANALYZE, BUFFERS)` sur un jeu de données simulé reproduisant la volumétrie attendue à 15 000 utilisateurs actifs, à l'aide du script `scripts/ml/benchmark_perf.py`. Les gains restent cohérents avec les ordres de grandeur habituellement observés lors du passage d'un `Seq Scan`/`Nested Loop` à un `Index Scan` sur PostgreSQL.*
