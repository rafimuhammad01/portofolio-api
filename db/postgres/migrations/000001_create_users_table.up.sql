CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    username VARCHAR (128) UNIQUE NOT NULL,
    password VARCHAR (64) NOT NULL,
    full_name VARCHAR (128) NOT NULL
);