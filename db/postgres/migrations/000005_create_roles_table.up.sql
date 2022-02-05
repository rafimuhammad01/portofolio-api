CREATE TABLE IF NOT EXISTS roles(
    id serial PRIMARY KEY,
    name VARCHAR (128) NOT NULL,
    description TEXT
);