-- =========================
-- USERS
-- =========================
INSERT INTO users (first_name, last_name, email, password_hash, role, score, created_at) VALUES
('Alice', 'Martin', 'alice@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'client', 10, NOW()),
('Bob', 'Durand', 'bob@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'pro', 50, NOW()),
('Charlie', 'Lefevre', 'charlie@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'admin', 100, NOW()),
('Diane', 'Moreau', 'diane@test.com', '$2b$10$hashhashhashhashhashhashhashhashhashhash', 'intra', 20, NOW());

-- =========================
-- EVENTS
-- =========================
INSERT INTO event (approved, approved_by, approved_at, price, date, created_by, created_at) VALUES
(true, 3, NOW(), 29.99, NOW() + INTERVAL '7 days', 2, NOW()),
(false, NULL, NULL, 0.00, NOW() + INTERVAL '14 days', 1, NOW());

-- Participation events
INSERT INTO event_participation (event_id, user_id) VALUES
(1, 1),
(1, 2);

-- =========================
-- THREADS & POSTS
-- =========================
INSERT INTO thread (created_by, title, content, upvotes, downvotes, created_at, last_post_at) VALUES
(1, 'Bienvenue', 'Premier thread de discussion', 5, 0, NOW(), NOW());

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
(2, NOW(), CURRENT_DATE + 3, '10:00', '12:00');

INSERT INTO entry_participation (entry_id, user_id, status, joined_at) VALUES
(1, 1, 'accepted', NOW()),
(1, 4, 'pending', NOW());

-- =========================
-- COURSES & ORDERS
-- =========================
INSERT INTO course (created_by, created_at, approved, approved_by, approved_at, price) VALUES
(2, NOW(), true, 3, NOW(), 99.99),
(2, NOW(), false, NULL, NULL, 49.99);

INSERT INTO course_order (course_id, buyer_id, price, booked_at) VALUES
(1, 1, 99.99, NOW());

-- =========================
-- CONTRACTS
-- =========================
INSERT INTO contract (name, created_by, content, created_at, until) VALUES
('Contrat coaching', 2, 'Contenu du contrat...', NOW(), CURRENT_DATE + 30);

-- =========================
-- LISTINGS & ORDERS
-- =========================
INSERT INTO listing (created_by, created_at, approved, approved_by, approved_at, status, price) VALUES
(2, NOW(), true, 3, NOW(), 'active', 150.00),
(2, NOW(), true, 3, NOW(), 'sold', 200.00);

INSERT INTO listing_order (listing_id, user_id, price, created_at, status) VALUES
(1, 1, 150.00, NOW(), 'paid'),
(2, 4, 200.00, NOW(), 'completed');
