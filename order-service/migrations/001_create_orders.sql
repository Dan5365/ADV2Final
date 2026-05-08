-- +goose Up
CREATE TABLE orders (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    status     VARCHAR(50) NOT NULL DEFAULT 'pending',
    total      NUMERIC(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id       BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    name     VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price    NUMERIC(12,2) NOT NULL
);

-- +goose Down
DROP TABLE order_items;
DROP TABLE orders;
