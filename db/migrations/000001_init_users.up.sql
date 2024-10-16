CREATE TABLE IF NOT EXISTS users
(
    id       bigserial PRIMARY KEY,
    username VARCHAR(50) UNIQUE  NOT NULL,
    email    VARCHAR(300) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL,
    salt     varchar(50)
);