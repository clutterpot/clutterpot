CREATE TABLE IF NOT EXISTS files(
    id char(20) PRIMARY KEY NOT NULL,
    name varchar(255) NOT NULL,
    mime_type varchar(127) NOT NULL,
    extension varchar(64) NOT NULL,
    size bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);
