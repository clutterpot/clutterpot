CREATE TABLE IF NOT EXISTS users(
    id char(20) PRIMARY KEY NOT NULL,
    username varchar(32) UNIQUE NOT NULL,
    password char(60) NOT NULL,
    email varchar(254) UNIQUE NOT NULL,
    kind smallint NOT NULL,
    display_name varchar(32),
    bio varchar(160),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
