-- 1. Создаем таблицу orders с внешними ключами
CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(100) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(100) NOT NULL,
    shardkey VARCHAR(50) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(50) NOT NULL
);


-- 2. Сначала создаем таблицу deliveries
CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    deliveries_id VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(50) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    
    CONSTRAINT deliveries_id_fk FOREIGN KEY (deliveries_id) REFERENCES orders (order_uid)
);


-- 3. Создаем таблицу payments
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(255),
	transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(100) NOT NULL,
    delivery_cost INTEGER NOT NULL CHECK (delivery_cost >= 0),
    goods_total INTEGER NOT NULL CHECK (goods_total >= 0),
    custom_fee INTEGER NOT NULL CHECK (custom_fee >= 0),

    CONSTRAINT transaction_id_fk FOREIGN KEY (transaction_id) REFERENCES orders (order_uid)
);

-- 4. Создаем таблицу items
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    items_id VARCHAR(255),
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    sale BIGINT NULL CHECK (sale >= 0),
    size VARCHAR(50) NOT NULL,
    total_price BIGINT NOT NULL CHECK (total_price >= 0),
    nm_id BIGINT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status BIGINT NOT NULL,

    CONSTRAINT items_id_fk FOREIGN KEY (items_id) REFERENCES orders (order_uid)
);


-- 5. Создаем индексы для улучшения производительности
CREATE INDEX idx_orders_order_uid ON orders(order_uid);
CREATE INDEX idx_orders_date_created ON orders(date_created);
CREATE INDEX idx_items_items_id ON items(items_id);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_deliveries_id ON deliveries(deliveries_id);