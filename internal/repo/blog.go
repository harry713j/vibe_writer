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

type authorData struct {
	FullName string `json:"author_name"`
	Bio      string `json:"author_bio"`
	Avatar   string `json:"author_avatar"`
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

func (b *BlogRepository) UpdateBlogVisibility(userId uuid.UUID, slug string) error {
	if _, err := b.DB.Exec("UPDATE blogs SET visibility = NOT visibility WHERE user_id = $1 AND slug = $2", userId, slug); err != nil {
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

// get all the public blogs of an user
func (b *BlogRepository) GetAllPublicBlog(userId uuid.UUID, page, limit int) (*model.PaginatedResponse[model.BlogSummary], error) {
	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var blogs []model.BlogSummary
	query := `
		SELECT 
			b.id,
			b.title,
			b.slug,
			b.content,
			b.visibility,
			b.user_id,
			b.created_at,
			b.updated_at,
			COALESCE((
				SELECT bp.photo_url 
				FROM blog_photos bp 
				WHERE bp.blog_id = b.id 
				ORDER BY bp.id ASC 
				LIMIT 1
			), '') AS blog_thumbnail,
			COUNT(l.*) FILTER (WHERE l.like_type = 'like')    AS likes_count,
			COUNT(l.*) FILTER (WHERE l.like_type = 'dislike') AS dislikes_count,
			COUNT(DISTINCT c.id) AS comments_count
		FROM blogs b
		LEFT JOIN likes l ON l.blog_id = b.id
		LEFT JOIN comments c ON c.blog_id = b.id
		WHERE b.user_id = $1 AND b.visibility = true
		GROUP BY b.id, b.title, b.slug, b.content, b.user_id, b.created_at, b.updated_at
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
		`

	rows, err := b.DB.Query(query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		blog := model.BlogSummary{}

		err := rows.Scan(
			&blog.Id,
			&blog.Title,
			&blog.UserId,
			&blog.Slug,
			&blog.Content,
			&blog.Visibility,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.Thumbnail,
			&blog.LikesCount,
			&blog.DislikeCount,
			&blog.CommentCount,
		)

		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	// total blogs
	var total int
	err = b.DB.QueryRow("SELECT COUNT(*) FROM blogs WHERE blogs.user_id = $1", userId).Scan(&total)

	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse[model.BlogSummary]{
		Data: blogs,
		Meta: model.PageMeta{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: totalPages,
		},
	}, nil
}

// get all the blogs of an user
func (b *BlogRepository) GetAllBlog(userId uuid.UUID, page, limit int) (*model.PaginatedResponse[model.BlogSummary], error) {
	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var blogs []model.BlogSummary
	query := `
		SELECT 
			b.id,
			b.title,
			b.slug,
			b.content,
			b.visibility,
			b.user_id,
			b.created_at,
			b.updated_at,
			COALESCE((
				SELECT bp.photo_url 
				FROM blog_photos bp 
				WHERE bp.blog_id = b.id 
				ORDER BY bp.id ASC 
				LIMIT 1
			), '') AS blog_thumbnail,
			COUNT(l.*) FILTER (WHERE l.like_type = 'like')    AS likes_count,
			COUNT(l.*) FILTER (WHERE l.like_type = 'dislike') AS dislikes_count,
			COUNT(DISTINCT c.id) AS comments_count
		FROM blogs b
		LEFT JOIN likes l ON l.blog_id = b.id
		LEFT JOIN comments c ON c.blog_id = b.id
		WHERE b.user_id = $1
		GROUP BY b.id, b.title, b.slug, b.content, b.user_id, b.created_at, b.updated_at
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
		`

	rows, err := b.DB.Query(query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		blog := model.BlogSummary{}

		err := rows.Scan(
			&blog.Id,
			&blog.Title,
			&blog.UserId,
			&blog.Slug,
			&blog.Content,
			&blog.Visibility,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.Thumbnail,
			&blog.LikesCount,
			&blog.DislikeCount,
			&blog.CommentCount,
		)

		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	// total blogs
	var total int
	err = b.DB.QueryRow("SELECT COUNT(*) FROM blogs WHERE blogs.user_id = $1", userId).Scan(&total)

	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse[model.BlogSummary]{
		Data: blogs,
		Meta: model.PageMeta{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: totalPages,
		},
	}, nil
}

// get blog by slug
func (b *BlogRepository) GetBlogBySlug(userId uuid.UUID, slug string) (*model.BlogResponse, error) {
	blogDataStat, err := b.getBlogWithStatBySlug(userId, slug)

	if err != nil {
		return nil, err
	}

	// author details
	authorData, err := b.getAuthorData(userId)

	if err != nil {
		return nil, err
	}

	comments, err := b.getBlogComments(blogDataStat.Id)

	if err != nil {
		return nil, err
	}

	blog := model.BlogResponse{
		BlogWithStat: *blogDataStat,
		Comments:     comments,
		AuthorName:   authorData.FullName,
		AuthorBio:    authorData.Bio,
		AuthorAvatar: authorData.Avatar,
	}

	return &blog, nil
}

// get blog by title
func (b *BlogRepository) GetBlogByTitle(userId uuid.UUID, title string) error {
	query := `SELECT * FROM blogs  WHERE user_id=$1 AND title=$2`
	_, err := b.DB.Exec(query, userId, title)

	if err != nil {
		return err
	}

	return nil
}

// get blog by id
func (b *BlogRepository) GetBlogById(userId uuid.UUID, blogId int64) (*model.BlogResponse, error) {
	var blog model.Blog

	err := b.DB.QueryRow("SELECT * FROM blogs WHERE id = $1", blogId).Scan(
		&blog.Id, &blog.UserId, &blog.Title, &blog.Slug, &blog.Content, &blog.Visibility,
		&blog.CreatedAt, &blog.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	blogData, err := b.getBlogWithStatBySlug(userId, blog.Slug)

	if err != nil {
		return nil, err
	}

	authorData, err := b.getAuthorData(userId)

	if err != nil {
		return nil, err
	}

	comments, err := b.getBlogComments(blogData.Id)

	if err != nil {
		return nil, err
	}

	blogRes := model.BlogResponse{
		BlogWithStat: *blogData,
		Comments:     comments,
		AuthorName:   authorData.FullName,
		AuthorBio:    authorData.Bio,
		AuthorAvatar: authorData.Avatar,
	}

	return &blogRes, nil
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

func (b *BlogRepository) getBlogWithStatBySlug(userId uuid.UUID, slug string) (*model.BlogWithStat, error) {
	var blogData model.BlogWithStat

	blogQuery := `
		SELECT 
			b.id,
			b.title, 
			b.user_id,
			b.slug,
			b.content,
			b.visibility,
			b.created_at,
			b.updated_at, 
			COALESCE(array_agg(bp.photo_url) FILTER (WHERE bp.photo_url IS NOT NULL), '{}') AS photo_urls,
			COUNT(l.*) FILTER (WHERE l.like_type = 'like') AS likes_count,
			COUNT(l.*) FILTER (WHERE l.like_type = 'dislike') AS dislikes_count

		FROM blogs b
		LEFT JOIN blog_photos bp ON b.id = bp.blog_id
		LEFT JOIN likes l ON b.id = l.blog_id
		WHERE b.user_id = $1 AND b.slug = $2
		GROUP BY b.id
	 `

	err := b.DB.QueryRow(blogQuery, userId, slug).Scan(
		&blogData.Id, &blogData.Title, &blogData.UserId, &blogData.Slug, &blogData.Content, &blogData.Visibility,
		&blogData.CreatedAt, &blogData.UpdatedAt, &blogData.PhotoUrls, &blogData.LikeCount, &blogData.DislikeCount,
	)

	if err != nil {
		return nil, err
	}

	return &blogData, nil
}

func (b *BlogRepository) getAuthorData(userId uuid.UUID) (*authorData, error) {
	var authorData authorData

	authorQuery := `
		SELECT full_name, bio, avatar_url FROM user_profiles WHERE user_id = $1
	`

	err := b.DB.QueryRow(authorQuery, userId).Scan(
		&authorData.FullName, &authorData.Bio, &authorData.Avatar,
	)

	if err != nil {
		return nil, err
	}

	return &authorData, nil
}

func (b *BlogRepository) getBlogComments(blogId int64) ([]model.CommentWithStat, error) {
	var comments []model.CommentWithStat

	commentQuery := `
		SELECT 
			c.id,
			c.content,
			c.user_id,
			c.parent_id,
			c.created_at,
			c.updated_at,
			COUNT(l.*) FILTER (WHERE l.like_type = 'like') AS likes_count,
			COUNT(l.*) FILTER (WHERE l.like_type = 'dislike') AS dislikes_count

		FROM comments c
		LEFT JOIN likes l ON c.id = l.comment_id
		WHERE c.blog_id = $1
		GROUP BY c.id
		ORDER BY c.created_at ASC
	`

	commentRows, err := b.DB.Query(commentQuery, blogId)

	if err != nil {
		return nil, err
	}

	defer commentRows.Close()

	for commentRows.Next() {
		var comment model.CommentWithStat

		err := commentRows.Scan(
			&comment.Id, &comment.Content, &comment.UserId, &comment.ParentId,
			&comment.CreatedAt, &comment.UpdatedAt, &comment.LikeCount, &comment.DislikeCount,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, err
}
