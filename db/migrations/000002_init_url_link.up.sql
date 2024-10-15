CREATE TABLE IF NOT EXISTS links
(
    id       bigserial PRIMARY KEY,
    url      TEXT,
    alias    VARCHAR(50) UNIQUE NOT NULL,
    users_id bigint REFERENCES users(id) NOT NULL
);