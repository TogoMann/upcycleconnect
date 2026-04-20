DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE USER_ROLE AS ENUM ('client', 'pro', 'interne', 'admin');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entry_status') THEN
        CREATE TYPE ENTRY_STATUS AS ENUM ('accepted', 'declined', 'pending');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'listing_categories') THEN
        CREATE TYPE LISTING_CATEGORIES AS ENUM ('Mobilier', 'Décoration', 'Vêtements', 'Jouet', 'Electronique', 'Outils');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'listing_status') THEN
        CREATE TYPE LISTING_STATUS AS ENUM ('active', 'sold', 'cancelled');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'listing_order_status') THEN
        CREATE TYPE LISTING_ORDER_STATUS AS ENUM ('pending', 'paid', 'shipped', 'completed');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'item_state') THEN
        CREATE TYPE ITEM_STATE AS ENUM ('neuf', 'bon etat', 'abime', 'casse');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'item_status') THEN
        CREATE TYPE ITEM_STATUS AS ENUM ('deposited', 'validated', 'collected');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'container_status') THEN
        CREATE TYPE CONTAINER_STATUS AS ENUM ('Available', 'Occupied', 'HS');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'container_size') THEN
        CREATE TYPE CONTAINER_SIZE AS ENUM ('S', 'M', 'L');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'thread_categories') THEN
        CREATE TYPE THREAD_CATEGORIES AS ENUM ('Bricolage', 'Textile', 'Ressources', 'Débutants', 'Communauté');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'project_status') THEN
        CREATE TYPE PROJECT_STATUS AS ENUM ('in progress', 'done', 'featured', 'cancelled');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ad_target_type') THEN
        CREATE TYPE AD_TARGET_TYPE AS ENUM ('listing', 'project', 'brand');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ad_type') THEN
        CREATE TYPE AD_TYPE AS ENUM ('partnership', 'other');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ad_status') THEN
        CREATE TYPE AD_STATUS AS ENUM ('pending', 'validated', 'rejected', 'expired');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS city (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    zip_code VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS address (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT REFERENCES city(id),
    street_name VARCHAR(255) NOT NULL,
    street_number VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS site (
    id BIGSERIAL PRIMARY KEY,
    address_id BIGINT REFERENCES address(id),
    type_site VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) UNIQUE NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    billing_cycle VARCHAR(20) DEFAULT 'monthly',
    features TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS container (
    id BIGSERIAL PRIMARY KEY,
    site_id BIGINT REFERENCES site(id),
    status CONTAINER_STATUS DEFAULT 'Available',
    size CONTAINER_SIZE DEFAULT 'M',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(64) UNIQUE NOT NULL,
    password_hash CHAR(60) NOT NULL,
    role USER_ROLE NOT NULL DEFAULT 'client',
    language_preference VARCHAR(5) DEFAULT 'fr',
    has_seen_tutorial BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS score_history (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    points INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event (
    id BIGSERIAL PRIMARY KEY,
    approved BOOLEAN NOT NULL DEFAULT FALSE,
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    price DECIMAL(6,2) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_participation (
    event_id BIGINT REFERENCES event(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    stripe_payment_intent_id VARCHAR(255),
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS thread (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    category THREAD_CATEGORIES NOT NULL,
    title VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_post_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS thread_views (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    thread_id BIGINT REFERENCES thread(id) ON  DELETE CASCADE,
    times INTEGER
);

CREATE TABLE IF NOT EXISTS post (
    id BIGSERIAL PRIMARY KEY,
    thread_id BIGINT NOT NULL REFERENCES thread(id) ON DELETE CASCADE,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    parent_id BIGINT REFERENCES post(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    news_id BIGINT NOT NULL REFERENCES news(id) ON DELETE CASCADE,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS entry (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    schedule DATE,
    start TIME,
    ending TIME
);

CREATE TABLE IF NOT EXISTS entry_participation (
    entry_id BIGINT REFERENCES entry(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status ENTRY_STATUS NOT NULL DEFAULT 'pending',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (entry_id, user_id)
);

CREATE TABLE IF NOT EXISTS course (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    max_capacity INTEGER,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL DEFAULT FALSE,
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    price DECIMAL(6,2)
);

CREATE TABLE IF NOT EXISTS course_order (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    buyer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stripe_payment_intent_id VARCHAR(255),
    price DECIMAL(6,2),
    booked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    subscriber_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL,
    tier VARCHAR(64) NOT NULL,
    created_at DATE DEFAULT CURRENT_DATE,
    until DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS item (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    container_id BIGINT REFERENCES container(id) ON DELETE SET NULL,
    site_id BIGINT REFERENCES site(id) ON DELETE SET NULL,
    material_type VARCHAR(64),
    physical_state ITEM_STATE,
    status ITEM_STATUS DEFAULT 'deposited',
    weight DECIMAL(8,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS container_access (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT REFERENCES item(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    access_code VARCHAR(16) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS listing (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    category LISTING_CATEGORIES,
    item_id BIGINT REFERENCES item(id) ON DELETE CASCADE,
    city_id BIGINT REFERENCES city(id) ON DELETE SET NULL,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL DEFAULT FALSE,
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    status LISTING_STATUS NOT NULL DEFAULT 'active',
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS listing_order (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listing(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stripe_payment_intent_id VARCHAR(255),
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status LISTING_ORDER_STATUS NOT NULL DEFAULT 'pending'
);

CREATE TABLE IF NOT EXISTS project (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT REFERENCES listing(id) ON DELETE SET NULL,
    creator_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(128) NOT NULL,
    description TEXT,
    final_score INTEGER,
    status PROJECT_STATUS DEFAULT 'in progress',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS project_steps (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT REFERENCES project(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS advertisement (
    id BIGSERIAL PRIMARY KEY,
    announcer_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    target_id BIGINT NOT NULL,
    target_type AD_TARGET_TYPE NOT NULL,
    ad_type AD_TYPE NOT NULL,
    budget DECIMAL(10,2),
    status AD_STATUS DEFAULT 'pending',
    stripe_payment_intent_id VARCHAR(255),
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL
);
