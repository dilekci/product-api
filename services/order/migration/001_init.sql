CREATE TABLE IF NOT EXISTS orders (
    id              BIGSERIAL PRIMARY KEY,
    customer_number VARCHAR(255) NOT NULL,
    product_id      VARCHAR(255) NOT NULL,
    quantity        INT          NOT NULL CHECK (quantity > 0),
    order_time      TIMESTAMP    NOT NULL DEFAULT NOW()
);