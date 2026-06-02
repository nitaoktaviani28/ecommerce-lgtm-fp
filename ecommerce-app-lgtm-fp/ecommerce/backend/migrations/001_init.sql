-- Migration: 001_init.sql
-- Inisialisasi schema database ecommerce

-- ─── Users ───────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ─── Products ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    description TEXT,
    price       DECIMAL(12, 2) NOT NULL CHECK (price >= 0),
    stock       INT            NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category    VARCHAR(100)   NOT NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_category ON products (category);

-- ─── Orders ───────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS orders (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT         NOT NULL REFERENCES users (id),
    status      VARCHAR(50)    NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending','paid','shipped','completed','cancelled')),
    total_price DECIMAL(12, 2) NOT NULL CHECK (total_price >= 0),
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status  ON orders (status);

-- ─── Order Items ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS order_items (
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT         NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id BIGINT         NOT NULL REFERENCES products (id),
    quantity   INT            NOT NULL CHECK (quantity > 0),
    price      DECIMAL(12, 2) NOT NULL CHECK (price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);

-- ─── Seed Data ────────────────────────────────────────────────────────────────
INSERT INTO users (name, email) VALUES
    ('Budi Santoso',  'budi@example.com'),
    ('Sari Dewi',     'sari@example.com'),
    ('Andi Pratama',  'andi@example.com')
ON CONFLICT (email) DO NOTHING;

INSERT INTO products (name, description, price, stock, category) VALUES
    ('Laptop Gaming ASUS ROG',    'Laptop gaming high-end dengan RTX 4060',       15000000, 10, 'electronics'),
    ('Smartphone Samsung Galaxy', 'Smartphone flagship dengan kamera 200MP',       8500000,  25, 'electronics'),
    ('Headphone Sony WH-1000XM5', 'Noise cancelling wireless headphone',           4200000,  30, 'electronics'),
    ('Kemeja Batik Premium',      'Batik tulis motif parang, bahan katun ATBM',     350000,  50, 'fashion'),
    ('Sepatu Sneakers Lokal',     'Sneakers handmade kulit sapi asli',              850000,  40, 'fashion'),
    ('Buku Clean Code',           'Panduan menulis kode yang bersih dan maintainable', 120000, 100, 'books'),
    ('Buku The Pragmatic Programmer', 'Buku klasik untuk software engineer',        150000,  80, 'books'),
    ('Kopi Toraja Arabika 500gr', 'Kopi single origin dari Toraja, Sulawesi',        95000, 200, 'food'),
    ('Sambal Oelek Homemade',     'Sambal oelek pedas manis buatan rumah',           35000, 150, 'food'),
    ('Tas Ransel Eiger',          'Tas ransel outdoor waterproof 30L',               750000,  35, 'outdoor')
ON CONFLICT DO NOTHING;
