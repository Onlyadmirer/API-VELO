-- VELO initial schema (PostgreSQL)
-- Derived from the queries used in pkg/repository.
-- Run: psql "$DATABASE_URL" -f migrations/001_init.sql

CREATE TABLE IF NOT EXISTS users (
    id                      SERIAL PRIMARY KEY,
    name                    VARCHAR(255) NOT NULL,
    email                   VARCHAR(255) NOT NULL UNIQUE,
    password                VARCHAR(255),                       -- NULL for OAuth users
    role                    VARCHAR(50)  NOT NULL DEFAULT 'customer',
    is_verified             BOOLEAN      NOT NULL DEFAULT FALSE,
    verify_token            VARCHAR(255),
    verify_token_expires_at TIMESTAMPTZ,
    oauth_id                VARCHAR(255),
    oauth_provider          VARCHAR(50),
    avatar_url              TEXT,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (oauth_provider, oauth_id)
);

CREATE INDEX IF NOT EXISTS idx_users_verify_token ON users (verify_token);

CREATE TABLE IF NOT EXISTS products (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255)   NOT NULL,
    stock      INT            NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category   VARCHAR(100),
    price      NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    image      TEXT,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS carts (
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cart_items (
    id         SERIAL PRIMARY KEY,
    cart_id    INT         NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
    product_id INT         NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    quantity   INT         NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cart_id, product_id)
);

CREATE TABLE IF NOT EXISTS orders (
    id           SERIAL PRIMARY KEY,
    user_id      INT            NOT NULL REFERENCES users (id),
    total_amount NUMERIC(12, 2) NOT NULL,
    status       VARCHAR(20)    NOT NULL DEFAULT 'Unpaid',  -- Unpaid | Pending | Paid | cancel
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status_created_at ON orders (status, created_at);

CREATE TABLE IF NOT EXISTS order_items (
    id         SERIAL PRIMARY KEY,
    order_id   INT            NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id INT            NOT NULL REFERENCES products (id),
    quantity   INT            NOT NULL CHECK (quantity > 0),
    price      NUMERIC(12, 2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);
