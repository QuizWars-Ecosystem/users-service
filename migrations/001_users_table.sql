-- Write your migrate up statements here

CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(32) UNIQUE NOT NULL,
    email VARCHAR(128) UNIQUE NOT NULL,
    pass_hash VARCHAR(64) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    avatar_id SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_users_name ON users(username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_users_email ON users(email);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_unique_users_email;
DROP INDEX IF EXISTS idx_unique_users_name;

DROP TABLE IF EXISTS users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
