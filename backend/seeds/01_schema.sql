DO $$
BEGIN
    BEGIN
        CREATE TYPE USER_ROLE AS ENUM ('client', 'pro', 'intra', 'admin');
    EXCEPTION
        WHEN duplicate_object THEN NULL;
    END;

    BEGIN
        CREATE TYPE ENTRY_STATUS AS ENUM ('accepted', 'declined', 'pending');
    EXCEPTION
        WHEN duplicate_object THEN NULL;
    END;

    BEGIN
        CREATE TYPE LISTING_STATUS AS ENUM ('active', 'sold', 'cancelled');
    EXCEPTION
        WHEN duplicate_object THEN NULL;
    END;

    BEGIN
        CREATE TYPE LISTING_ORDER_STATUS AS ENUM ('pending', 'paid', 'shipped', 'completed');
    EXCEPTION
        WHEN duplicate_object THEN NULL;
    END;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL,
    password_hash CHAR(60) NOT NULL,
    role USER_ROLE NOT NULL,
    score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event (
    id BIGSERIAL PRIMARY KEY,
    approved BOOLEAN NOT NULL,
    approved_by BIGINT REFERENCES users(id),
    approved_at TIMESTAMP,
    price DECIMAL(6,2) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_participation (
    event_id BIGINT REFERENCES event(id),
    user_id BIGINT REFERENCES users(id),
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS thread (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(64) NOT NULL,
    content VARCHAR(255) NOT NULL,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_post_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post (
    id BIGSERIAL PRIMARY KEY,
    thread_id BIGINT NOT NULL REFERENCES thread(id),
    created_by BIGINT NOT NULL REFERENCES users(id),
    content VARCHAR(255) NOT NULL,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(64) NOT NULL,
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    news_id BIGINT NOT NULL REFERENCES news(id),
    created_by BIGINT NOT NULL REFERENCES users(id),
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS entry (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    schedule DATE,
    start TIME,
    ending TIME
);

CREATE TABLE IF NOT EXISTS entry_participation (
    entry_id BIGINT REFERENCES entry(id),
    user_id BIGINT REFERENCES users(id),
    PRIMARY KEY (entry_id, user_id),

    status ENTRY_STATUS NOT NULL DEFAULT 'pending',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS course (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL,
    approved_by BIGINT REFERENCES users(id) DEFAULT NULL,
    approved_at TIMESTAMP,
    price DECIMAL(6,2)
);

CREATE TABLE IF NOT EXISTS course_order (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES course(id),
    buyer_id BIGINT NOT NULL REFERENCES users(id),
    price DECIMAL(6,2),
    booked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contract (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    until DATE
);

CREATE TABLE IF NOT EXISTS listing (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL,
    approved_by BIGINT NOT NULL REFERENCES users(id),
    approved_at TIMESTAMP,
    status LISTING_STATUS NOT NULL DEFAULT 'active',
    price DECIMAL(6, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS listing_order (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listing(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    price DECIMAL(6, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status LISTING_ORDER_STATUS NOT NULL DEFAULT 'pending'
);
