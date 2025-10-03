-- +goose Up
CREATE TABLE IF NOT EXISTS blogs(
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    content TEXT NOT NULL,
    visibility BOOLEAN DEFAULT true,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_slug UNIQUE(user_id, slug)
);

CREATE TABLE IF NOT EXISTS blog_photos(
    id BIGSERIAL PRIMARY KEY,
    blog_id BIGINT NOT NULL,
    photo_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_blog
    FOREIGN KEY(blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS blog_photos;
DROP TABLE IF EXISTS blogs;
