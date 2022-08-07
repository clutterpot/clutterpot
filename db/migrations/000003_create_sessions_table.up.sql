CREATE TABLE IF NOT EXISTS sessions (
    id CHAR(20) PRIMARY KEY NOT NULL,
    user_id CHAR(20) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 year' NOT NULL
);
