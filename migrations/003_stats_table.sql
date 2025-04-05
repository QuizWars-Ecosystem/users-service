-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS stats (
    user_id VARCHAR(36) PRIMARY KEY REFERENCES users(id),
    rating SMALLINT NOT NULL DEFAULT 0,
    coins INTEGER NOT NULL DEFAULT 0
);

---- create above / drop below ----

DROP TABLE IF EXISTS stats;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
