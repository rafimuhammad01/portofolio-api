CREATE TABLE IF NOT EXISTS projects(
    id serial PRIMARY KEY,
    name VARCHAR (128) NOT NULL,
    client_name VARCHAR (128) NOT NULL,
    description TEXT
);