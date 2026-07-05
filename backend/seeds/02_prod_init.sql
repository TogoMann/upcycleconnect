-- =========================
-- SYSTEM REFERENCE DATA
-- =========================

-- Default Cities
INSERT INTO city (name, zip_code) VALUES
('Paris', '75000'),
('Lyon', '69000'),
('Marseille', '13000'),
('Nantes', '44000'),
('Bordeaux', '33000'),
('Lille', '59000'),
('Strasbourg', '67000'),
('Montpellier', '34000'),
('Toulouse', '31000'),
('Nice', '06000')
ON CONFLICT DO NOTHING;

-- Default Plans
INSERT INTO plans (name, description, price, billing_cycle, features, is_active) VALUES
('Free', 'Accès limité aux fonctionnalités de base', 0.00, 'monthly', '{"Vente d''objets", "Dépôt d''objets"}', true),
('Premium', 'Accès complet à toutes les fonctionnalités', 9.99, 'monthly', '{"Vente d''objets", "Dépôt d''objets", "Ateliers gratuits", "Événements VIP"}', true),
('Pro', 'Pour les professionnels de l''upcycling', 29.99, 'monthly', '{"Fonctionnalités Premium", "Statistiques avancées", "Publicité prioritaire"}', true)
ON CONFLICT (name) DO UPDATE 
SET description = EXCLUDED.description,
    price = EXCLUDED.price,
    billing_cycle = EXCLUDED.billing_cycle,
    features = EXCLUDED.features,
    is_active = EXCLUDED.is_active;

-- =========================
-- DEFAULT ADMIN USER
-- =========================
INSERT INTO users (username, first_name, last_name, email, password_hash, role, created_at) VALUES
('clefevre', 'Charlie', 'Lefevre', 'charlie@test.com', '$2a$10$Lp3TJGSMCLYl1sFOr14C/ummqFAS6avxKrd5mjxJ0qjS.gcb5VzIa', 'admin', NOW())
ON CONFLICT (username) DO UPDATE
SET password_hash = EXCLUDED.password_hash;
