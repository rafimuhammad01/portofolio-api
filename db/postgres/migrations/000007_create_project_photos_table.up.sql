CREATE TABLE IF NOT EXISTS project_photos(
    id serial PRIMARY KEY,
    photo TEXT,
    description TEXT,
    project_id INTEGER REFERENCES projects(id)
);