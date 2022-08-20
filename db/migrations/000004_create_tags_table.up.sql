CREATE TABLE IF NOT EXISTS tags(
    id char(20) PRIMARY KEY NOT NULL,
    owner_id char(20) REFERENCES users(id) ON DELETE CASCADE DEFAULT NULL,
    name varchar(32) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone,
    UNIQUE NULLS NOT DISTINCT(owner_id, name)
);
