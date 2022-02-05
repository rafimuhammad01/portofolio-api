CREATE TABLE IF NOT EXISTS jastip_members(
    id serial PRIMARY KEY,
    name VARCHAR (128) NOT NULL,
    photo TEXT
);