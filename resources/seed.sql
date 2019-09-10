-- Db
CREATE DATABASE reddit;

-- Admin service account
CREATE USER admin WITH ENCRYPTED PASSWORD 'admin';
GRANT ALL PRIVILEGES ON DATABASE reddit TO admin;

-- Reddit Posts
CREATE TABLE IF NOT EXISTS posts (
    id serial primary key,
    url varchar(500),
    text text
)

-- Admin service account
GRANT ALL PRIVILEGES ON TABLE posts TO admin;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO admin;
