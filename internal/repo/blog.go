package repo

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
)

type BlogRepository struct {
	DB *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{
		DB: db,
	}
}

// create blog
func (b *BlogRepository) CreateBlog(userId uuid.UUID, title, slug, content string) (int64, error) {
	blog := &model.Blog{
		UserId:    userId,
		Title:     title,
		Slug:      slug,
		Content:   content,
		CreatedAt: time.Now(),
	}

	var blogId int64

	err := b.DB.QueryRow(`INSERT INTO blogs(user_id, title, slug, content, created_at) 
		VALUES($1, $2, $3, $4, $5)
		RETURNING id`,
		blog.UserId, blog.Title, blog.Slug, blog.Content, blog.CreatedAt).Scan(&blogId)

	if err != nil {
		return 0, err
	}

	return blogId, nil
}

// create blog image
func (b *BlogRepository) CreateBlogImage(blogId int64, photoUrl string) (int64, error) {
	var blogPhotoId int64

	err := b.DB.QueryRow("INSERT INTO blog_photos(blog_id, photo_url) VALUES($1, $2) RETURNING id",
		blogId, photoUrl).Scan(&blogPhotoId)

	if err != nil {
		return 0, err
	}

	return blogPhotoId, nil
}

// update by slug -> slug is immutable
func (b *BlogRepository) UpdateBlog(userId uuid.UUID, slug string, title string, content string) (int64, error) {
	var blogId int64

	err := b.DB.QueryRow(`UPDATE blogs SET title=$1, content=$2 WHERE slug=$3 AND user_id=$4 RETURNING id`,
		title, content, slug, userId).Scan(&blogId)

	if err != nil {
		return 0, err
	}

	return blogId, nil
}

// update blog photos
func (b *BlogRepository) UpdateBlogPhoto(blogId, blogPhotoId int64, photoUrl string) error {
	if _, err := b.DB.Exec("UPDATE blog_photos SET photo_url=$1 WHERE blog_id=$2 AND id=$3",
		photoUrl, blogId, blogPhotoId); err != nil {
		return err
	}

	return nil
}

// delete by slug
func (b *BlogRepository) DeleteBlog(userId uuid.UUID, slug string) error {
	if _, err := b.DB.Exec("DELETE FROM blogs WHERE slug=$1 AND user_id=$2", slug, userId); err != nil {
		return err
	}
	return nil
}

// delete blog photos of a blog
func (b *BlogRepository) DeleteBlogPhotosByURLs(blogId int64, photoURLs []string) error {
	_, err := b.DB.Exec(
		"DELETE FROM blog_photos WHERE blog_id=$1 AND photo_url = ANY($2)",
		blogId,
		photoURLs,
	)
	return err
}

// delete blog photo
func (b *BlogRepository) DeleteBlogPhoto(blogId int64, photoUrl string) error {
	if _, err := b.DB.Exec("DELETE FROM blog_photos WHERE photo_url=$1 AND blog_id=$2", photoUrl, blogId); err != nil {
		return err
	}
	return nil
}

// get all the blogs of an user
func (b *BlogRepository) GetAllBlog(userId uuid.UUID) ([]*model.BlogDetails, error) {
	var blogs []*model.BlogDetails

	query := `SELECT b.id, b.title, b.user_id, b.slug, b.content, b.created_at, b.updated_at, 
	COALESCE(array_agg(bp.photo_url) FILTER (WHERE bp.photo_url IS NOT NULL), '{}')
		 FROM blogs b LEFT JOIN blog_photos bp ON b.id = bp.blog_id WHERE b.user_id=$1 GROUP BY b.id`

	rows, err := b.DB.Query(query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		blog := &model.BlogDetails{}

		err := rows.Scan(
			&blog.Id,
			&blog.Title,
			&blog.UserId,
			&blog.Slug,
			&blog.Content,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PhotoUrls,
		)

		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// get blog by slug
func (b *BlogRepository) GetBlogBySlug(userId uuid.UUID, slug string) (*model.BlogDetails, error) {
	var blog model.BlogDetails

	query := `SELECT b.id, b.title, b.user_id, b.slug, b.content, b.created_at, b.updated_at, 
	COALESCE(array_agg(bp.photo_url) FILTER (WHERE bp.photo_url IS NOT NULL), '{}')
		 FROM blogs b LEFT JOIN blog_photos bp ON b.id = bp.blog_id WHERE b.user_id=$1 AND b.slug=$2 GROUP BY b.id`

	err := b.DB.QueryRow(query, userId, slug).Scan(
		&blog.Id,
		&blog.Title,
		&blog.UserId,
		&blog.Slug,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.PhotoUrls,
	)

	if err != nil {
		return nil, err
	}

	return &blog, nil
}

// get blog by id
func (b *BlogRepository) GetBlogById(userId uuid.UUID, blogId int64) (*model.BlogDetails, error) {
	var blog model.BlogDetails

	query := `SELECT b.id, b.title, b.user_id, b.slug, b.content, b.created_at, b.updated_at, 
	COALESCE(array_agg(bp.photo_url) FILTER (WHERE bp.photo_url IS NOT NULL), '{}')
		 FROM blogs b LEFT JOIN blog_photos bp ON b.id = bp.blog_id WHERE b.user_id=$1 AND b.id=$2 GROUP BY b.id`

	err := b.DB.QueryRow(query, userId, blogId).Scan(
		&blog.Id,
		&blog.Title,
		&blog.UserId,
		&blog.Slug,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.PhotoUrls,
	)

	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (b *BlogRepository) GetPhotoUrls(blogId int64) ([]string, error) {
	var photoUrls []string
	rows, err := b.DB.Query("SELECT photo_url FROM blog_photos WHERE blog_id=$1", blogId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}

		photoUrls = append(photoUrls, url)
	}

	return photoUrls, nil
}
