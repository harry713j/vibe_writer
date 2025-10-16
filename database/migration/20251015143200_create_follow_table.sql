-- +goose Up
CREATE TABLE IF NOT EXISTS follows(
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_follow PRIMARY KEY(follower_id, following_id),
    CONSTRAINT fk_follower FOREIGN KEY(follower_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_following FOREIGN KEY(following_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_no_self_follow CHECK(follower_id <> following_id)
);

CREATE INDEX IF NOT EXISTS idx_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_following_id ON follows(following_id);

-- +goose Down
DROP INDEX IF EXISTS idx_following_id;
DROP INDEX IF EXISTS idx_follower_id;
DROP TABLE IF EXISTS follows;
