CREATE TABLE IF NOT EXISTS products (
    id          SERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    price       NUMERIC(10,2) NOT NULL,
    image_url   TEXT        NOT NULL DEFAULT '',
    stock       INTEGER     NOT NULL DEFAULT 0
);

CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipped', 'completed', 'cancelled');

CREATE TABLE IF NOT EXISTS orders (
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status     order_status NOT NULL DEFAULT 'pending',
    total      NUMERIC(10,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id         SERIAL PRIMARY KEY,
    order_id   INTEGER     NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER     NOT NULL REFERENCES products(id),
    name       TEXT        NOT NULL,
    quantity   INTEGER     NOT NULL,
    price      NUMERIC(10,2) NOT NULL
);
