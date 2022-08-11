CREATE TABLE IF NOT EXISTS tags(
    id char(20) PRIMARY KEY NOT NULL,
    owner_id char(20) REFERENCES users(id) ON DELETE CASCADE DEFAULT NULL,
    name varchar(32) NOT NULL,
    UNIQUE NULLS NOT DISTINCT(owner_id, name)
);
