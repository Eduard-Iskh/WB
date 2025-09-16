-- Создание базы данных
CREATE DATABASE orders_db;

-- Создание пользователя (альтернативный подход)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'orders_user') THEN
        CREATE USER orders_user WITH PASSWORD 'orders_password';
    END IF;
END
$$;

-- Предоставление привилегий на базу данных
GRANT ALL PRIVILEGES ON DATABASE orders_db TO orders_user;

-- Подключение к созданной базе данных
\c orders_db;

-- Создание таблицы в правильной базе данных
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR PRIMARY KEY,
    data JSONB NOT NULL
);

-- Предоставление привилегий на таблицу
GRANT ALL PRIVILEGES ON TABLE orders TO orders_user;