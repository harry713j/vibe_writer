-- +goose Up
CREATE TABLE IF NOT EXISTS comments(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    blog_id BIGINT NOT NULL,
    parent_id BIGINT, -- for nested comments
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog
    FOREIGN KEY(blog_id) REFERENCES blogs(id) ON DELETE CASCADE,
    CONSTRAINT fk_parent 
    FOREIGN KEY(parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

CREATE TYPE liketype AS ENUM('like', 'dislike');

CREATE TABLE IF NOT EXISTS likes(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    blog_id BIGINT,
    comment_id BIGINT,
    like_type LIKETYPE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_blog
    FOREIGN KEY(blog_id) REFERENCES blogs(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment
    FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    CONSTRAINT like_check CHECK(
        (blog_id IS NOT NULL AND comment_id IS NULL) OR
        (blog_id IS NULL AND comment_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX unique_user_blog_like 
ON likes(user_id, blog_id) 
WHERE blog_id IS NOT NULL;

CREATE UNIQUE INDEX unique_user_comment_like 
ON likes(user_id, comment_id) 
WHERE comment_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS unique_user_comment_like;
DROP INDEX IF EXISTS unique_user_blog_like;

DROP TABLE IF EXISTS likes;
DROP TYPE liketype;
DROP TABLE IF EXISTS comments;
