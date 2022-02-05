CREATE TABLE IF NOT EXISTS jastip_member_roles(
    id serial PRIMARY KEY,
    jastip_member_id INTEGER REFERENCES jastip_members(id),
    role_id INTEGER REFERENCES roles(id)
);