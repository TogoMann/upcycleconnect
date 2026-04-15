-- =========================
-- CITY & ADDRESS & SITE & CONTAINER
-- =========================
INSERT INTO city (name, zip_code) VALUES
('Paris', '75000'),
('Lyon', '69000'),
('Marseille', '13000');

INSERT INTO address (city_id, street_name, street_number) VALUES
(1, 'Rue de Rivoli', '1'),
(2, 'Rue de la République', '10'),
(3, 'La Canebière', '100');

INSERT INTO site (address_id, type_site, created_at) VALUES
(1, 'Point de collecte', NOW()),
(2, 'Atelier de réparation', NOW()),
(3, 'Centre de tri', NOW());

INSERT INTO container (site_id, status, size, created_at) VALUES
(1, 'Available', 'M', NOW()),
(1, 'Occupied', 'L', NOW()),
(2, 'HS', 'S', NOW());

-- =========================
-- USERS
-- =========================
INSERT INTO users (username, first_name, last_name, email, password_hash, role, created_at) VALUES
('amartin', 'Alice', 'Martin', 'alice@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'client', NOW()),
('bdurand', 'Bob', 'Durand', 'bob@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'pro', NOW()),
('clefevre', 'Charlie', 'Lefevre', 'charlie@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'admin', NOW()),
('dmoreau', 'Diane', 'Moreau', 'diane@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'interne', NOW()),
('mdede', 'morad', 'dede', 'moradtest@test.com', '$2a$12$7bg7UBVasAqV9aah61WcC.b25cw/lmKwR0dbJ/iVOuP1UDpIVmrOS', 'client', NOW());

-- =========================
-- SCORE HISTORY
-- =========================
INSERT INTO score_history (user_id, points, description, created_at) VALUES
(1, 10, 'Premier don d''objet', NOW()),
(1, 5, 'Commentaire utile', NOW()),
(2, 50, 'Vente finalisée', NOW());

-- =========================
-- ITEMS
-- =========================
INSERT INTO item (owner_id, container_id, site_id, material_type, physical_state, status, weight, created_at) VALUES
(1, 1, 1, 'Bois', 'bon etat', 'deposited', 5.5, NOW()),
(2, 2, 1, 'Metal', 'neuf', 'validated', 12.0, NOW());

-- =========================
-- EVENTS
-- =========================
INSERT INTO event (approved, approved_by, approved_at, price, date, created_by, created_at) VALUES
(true, 3, NOW(), 29.99, NOW() + INTERVAL '7 days', 2, NOW()),
(false, NULL, NULL, 0.00, NOW() + INTERVAL '14 days', 1, NOW());

-- Participation events
INSERT INTO event_participation (event_id, user_id, stripe_payment_intent_id) VALUES
(1, 1, 'pi_123456789'),
(1, 2, 'pi_987654321');

-- =========================
-- THREADS & POSTS
-- =========================
INSERT INTO thread (created_by, category, title, content, upvotes, downvotes, created_at, last_post_at) VALUES
(1, 'Communauté', 'Bienvenue', 'Premier thread de discussion', 5, 0, NOW(), NOW());

INSERT INTO post (thread_id, created_by, content, upvotes, downvotes, created_at) VALUES
(1, 2, 'Super idée !', 3, 0, NOW()),
(1, 3, 'Merci pour ce thread', 2, 0, NOW());

-- =========================
-- NEWS & COMMENTS
-- =========================
INSERT INTO news (created_by, title, content, created_at, upvotes, downvotes) VALUES
(3, 'Nouvelle fonctionnalité', 'On a ajouté plein de choses !', NOW(), 10, 1);

INSERT INTO comments (news_id, created_by, content, created_at, upvotes, downvotes) VALUES
(1, 1, 'Génial !', NOW(), 2, 0),
(1, 2, 'Hâte de tester', NOW(), 1, 0);

-- =========================
-- ENTRY & PARTICIPATION
-- =========================
INSERT INTO entry (created_by, created_at, schedule, start, ending) VALUES
(2, NOW(), CURRENT_DATE + 3, '10:00:00', '12:00:00');

INSERT INTO entry_participation (entry_id, user_id, status, joined_at) VALUES
(1, 1, 'accepted', NOW()),
(1, 4, 'pending', NOW());

-- =========================
-- COURSES & ORDERS
-- =========================
INSERT INTO course (name, description, max_capacity, created_by, created_at, approved, approved_by, approved_at, price) VALUES
('Couture Débutant', 'Apprenez les bases de la couture.', 10, 2, NOW(), true, 3, NOW(), 99.99),
('Menuiserie', 'Travail du bois pour tous.', 5, 2, NOW(), false, NULL, NULL, 49.99);

INSERT INTO course_order (course_id, buyer_id, price, booked_at, stripe_payment_intent_id) VALUES
(1, 1, 99.99, NOW(), 'pi_course_1');

-- =========================
-- SUBSCRIPTIONS
-- =========================
INSERT INTO subscriptions (subscriber_id, price, tier, created_at, until) VALUES
(2, 49.99, 'Premium', CURRENT_DATE, CURRENT_DATE + 30);

-- =========================
-- LISTINGS & ORDERS
-- =========================
INSERT INTO listing (name, description, category, item_id, city_id, created_by, created_at, approved, approved_by, approved_at, status, price) VALUES
('Chaise en bois', 'Une belle chaise faite main.', 'Mobilier', 1, 1, 2, NOW(), true, 3, NOW(), 'active', 150.00),
('Etagère métal', 'Etagère solide en métal.', 'Mobilier', 2, 2, 2, NOW(), true, 3, NOW(), 'sold', 200.00);

INSERT INTO listing_order (listing_id, user_id, price, created_at, status, stripe_payment_intent_id) VALUES
(1, 1, 150.00, NOW(), 'paid', 'pi_listing_1'),
(2, 4, 200.00, NOW(), 'completed', 'pi_listing_2');

-- =========================
-- PROJECTS & STEPS
-- =========================
INSERT INTO project (listing_id, creator_id, title, description, final_score, status, created_at) VALUES
(1, 1, 'Restauration Chaise', 'Ponçage et vernissage.', 80, 'in progress', NOW());

INSERT INTO project_steps (project_id, step_number, description, created_at) VALUES
(1, 1, 'Achat du matériel', NOW()),
(1, 2, 'Ponçage de la structure', NOW());
