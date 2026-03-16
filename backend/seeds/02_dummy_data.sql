-- 1. Insert Users
INSERT INTO users (first_name, last_name, email, password_hash, score, role) VALUES
('Alice', 'Admin', 'alice@system.com', '$2a$12$VQ7V59.LzUqZpA.o6.N5fO1J1V8P3zK6Z1/wVz5Rz/8xR9yO8W/mS', 1000, 'admin'),
('Bob', 'Pro', 'bob@work.com', '$2a$12$VQ7V59.LzUqZpA.o6.N5fO1J1V8P3zK6Z1/wVz5Rz/8xR9yO8W/mS', 500, 'pro'),
('Charlie', 'Client', 'charlie@web.com', '$2a$12$VQ7V59.LzUqZpA.o6.N5fO1J1V8P3zK6Z1/wVz5Rz/8xR9yO8W/mS', 50, 'client'),
('Dave', 'Intra', 'dave@corp.com', '$2a$12$VQ7V59.LzUqZpA.o6.N5fO1J1V8P3zK6Z1/wVz5Rz/8xR9yO8W/mS', 200, 'intra');

-- 2. Events
INSERT INTO event (approved, price, date, created_by, created_at) VALUES
(TRUE, 49.99, '2026-05-20 18:00:00', 1, NOW()),
(FALSE, 15.00, '2026-06-15 10:00:00', 2, NOW());

-- 3. Event Participation
INSERT INTO event_participation (event_id, user_id) VALUES
(1, 2),
(1, 3);

-- 4. Threads & Posts (Forum)
INSERT INTO thread (user_id, title, content, upvotes, downvotes, created_at) VALUES
(3, 'Welcome to the Forum', 'Happy to be here!', 10, 1, NOW() - INTERVAL '2 days');

INSERT INTO post (thread_id, user_id, content, created_at) VALUES
(1, 1, 'Welcome Charlie! Let us know if you need help.', NOW() - INTERVAL '1 day');

-- 5. News & Comments
INSERT INTO news (created_by, title, content, created_at, upvotes, downvotes) VALUES
(1, 'System Update v2.0', 'We have added new features to the dashboard.', NOW(), 50, 0);

INSERT INTO comments (user_id, news_id, content, created_at, upvotes, downvotes) VALUES
(2, 1, 'Great update, the UI feels much faster.', NOW(), 5, 0);

-- 6. Entries (Schedule/Logs)
INSERT INTO entry (created_by, created_at, date, start, end) VALUES
(2, NOW(), '2026-04-01', '09:00:00', '17:00:00');

INSERT INTO entry_participation (entry_id, user_id, status, joined_at) VALUES
(1, 3, 'accepted', NOW());

-- 7. Courses & Course Orders
INSERT INTO course (created_by, approved, approved_by, approved_at, price) VALUES
(2, 1, 1, NOW(), 199.99);

INSERT INTO course_order (course_id, buyer_id, price, booked_at) VALUES
(1, 3, 199.99, NOW());

-- 8. Contracts
INSERT INTO contract (name, created_by, content, created_at, until) VALUES
('Service Agreement', 1, 'Standard terms and conditions...', NOW(), '2027-12-31');

-- 9. Listings & Listing Orders
INSERT INTO listing (created_by, created_at, approved, approved_by, approved_at, status, price) VALUES
(2, NOW(), 1, 1, NOW(), 'active', 25.50);

INSERT INTO listing_order (listing_id, user_id, price, created_at, status) VALUES
(1, 4, 25.50, NOW(), 'paid');
