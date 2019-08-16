--DB Seed
CREATE DATABASE reddit;

--Admin Service Acct
CREATE USER admin WITH ENCRYPTED PASSWORD 'admin';
GRANT ALL PRIVILEGES ON DATABASE reddit TO admin;
GRANT ALL PRIVILEGES ON TABLE posts TO admin;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO admin;

--Reddit posts
CREATE TABLE posts (
    id serial primary key,
    url varchar(500),
    text text
)