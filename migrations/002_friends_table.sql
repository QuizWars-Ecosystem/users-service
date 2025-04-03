-- Write your migrate up statements here

CREATE TYPE friend_status AS ENUM ('pending', 'accepted', 'blocked');

CREATE TABLE IF NOT EXISTS friends (
    user_id UUID NOT NULL REFERENCES users(id),
    friend_id UUID NOT NULL REFERENCES users(id),
    status friend_status NOT NULL DEFAULT 'pending',
    PRIMARY KEY (user_id, friend_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS friends;

DROP TYPE IF EXISTS friend_status;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
