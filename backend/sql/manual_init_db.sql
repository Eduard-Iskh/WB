DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_catalog.pg_roles
        WHERE rolname = 'orders_user'
    ) THEN
        CREATE ROLE orders_user LOGIN PASSWORD 'orders_password';
    END IF;
END
$$;

CREATE DATABASE orders_db OWNER orders_user;

GRANT ALL PRIVILEGES ON DATABASE orders_db TO orders_user;