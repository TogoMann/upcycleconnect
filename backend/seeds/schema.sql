DO $$
BEGIN
    CREATE TYPE user_role AS ENUM ('client', 'pro', 'intra', 'admin');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL,
    password_hash CHAR(60) NOT NULL,
    role user_role NOT NULL
);

CREATE TABLE IF NOT EXISTS event (
    id BIGSERIAL PRIMARY KEY,
    approved BOOLEAN NOT NULL,
    price DECIMAL(6,2) NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS event_participation (
    event_id BIGINT REFERENCES event(id),
    user_id BIGINT REFERENCES users(id),
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS thread (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(64) NOT NULL,
    content VARCHAR(255) NOT NULL,
    upvotes INTEGER NOT NULL,
    downvotes INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_post_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post (
    id BIGSERIAL PRIMARY KEY,
    thread_id BIGINT NOT NULL REFERENCES thread(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    edited_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(64) NOT NULL,
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    upvotes INTEGER NOT NULL,
    downvotes INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    news_id BIGINT NOT NULL REFERENCES news(id),
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    upvotes INTEGER NOT NULL,
    downvotes INTEGER NOT NULL
);
