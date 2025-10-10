-- +goose Up
CREATE TABLE IF NOT EXISTS bookmarks(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    blog_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog FOREIGN KEY(blog_id) REFERENCES blogs(id) ON DELETE CASCADE,
    CONSTRAINT unique_bookmark UNIQUE(user_id, blog_id)
);

-- +goose Down
DROP TABLE IF EXISTS bookmarks;

