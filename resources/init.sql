-- Reddit Posts
CREATE TABLE IF NOT EXISTS posts (
    id serial primary key,
    url varchar(500),
    text text
);
