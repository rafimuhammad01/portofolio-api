CREATE TABLE IF NOT EXISTS skills(
    id serial PRIMARY KEY,
    skill VARCHAR (128) NOT NULL,
    description TEXT,
    jastip_member_id integer REFERENCES jastip_members (id)
);