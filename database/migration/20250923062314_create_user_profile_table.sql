-- +goose Up
CREATE TABLE IF NOT EXISTS user_profiles(
    user_id UUID PRIMARY KEY,
    full_name TEXT,
    bio TEXT,
    avatar_url TEXT,

    CONSTRAINT fk_user
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE IF EXISTS user_profiles;

