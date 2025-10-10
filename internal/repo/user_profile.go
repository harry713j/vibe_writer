package repo

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
)

type UserProfileRepository struct {
	DB *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *UserProfileRepository {
	return &UserProfileRepository{
		DB: db,
	}
}

// create profile
func (u *UserProfileRepository) CreateUserProfile(userId uuid.UUID) (*model.UserProfile, error) {
	profile := &model.UserProfile{
		UserId: userId,
	}

	if _, err := u.DB.Exec("INSERT INTO user_profiles(user_id) VALUES($1)", profile.UserId); err != nil {
		return nil, err
	}

	return profile, nil
}

func (u *UserProfileRepository) UpdateProfile(userId uuid.UUID, fullName, bio string) error {

	if _, err := u.DB.Exec("UPDATE user_profiles SET full_name=$1, bio=$2 WHERE user_id=$3",
		fullName, bio, userId); err != nil {
		return err
	}

	return nil
}

func (u *UserProfileRepository) UpdateAvatar(userId uuid.UUID, avatarUrl string) error {
	if _, err := u.DB.Exec("UPDATE user_profiles SET avatar_url=$1 WHERE user_id=$2", avatarUrl, userId); err != nil {
		return err
	}

	return nil
}

func (u *UserProfileRepository) GetUserDetails(userId uuid.UUID) (*model.UserDetails, error) {
	var userData model.UserDetails

	query := `SELECT u.id, u.username, u.email, u.created_at, u.updated_at, 
	COALESCE(up.full_name,'') AS full_name, COALESCE(up.bio, '') AS bio,
	COALESCE(up.avatar_url, '') AS avatar_url
	FROM user_profiles up INNER JOIN users u ON up.user_id = u.id 
	WHERE up.user_id = $1`

	err := u.DB.QueryRow(query, userId).Scan(
		&userData.Id,
		&userData.Username,
		&userData.Email,
		&userData.CreatedAt,
		&userData.UpdatedAt,
		&userData.FullName,
		&userData.Bio,
		&userData.AvatarUrl,
	)

	if err != nil {
		return nil, err
	}

	return &userData, nil
}

func (u *UserProfileRepository) GetAvatarUrl(userId uuid.UUID) (string, error) {
	var avatarurl string
	err := u.DB.QueryRow("SELECT COALESCE(avatar_url,'') AS avatar_url FROM user_profiles WHERE user_id=$1", userId).Scan(&avatarurl)

	if err != nil {
		return "", err
	}

	return avatarurl, nil
}

func (u *UserProfileRepository) DeleteAvatarUrl(userId uuid.UUID) error {
	if _, err := u.DB.Exec("UPDATE user_profiles SET avatar_url=NULL WHERE user_id=$1", userId); err != nil {
		return err
	}

	return nil
}

func (u *UserProfileRepository) GetAllBookmarks(userId uuid.UUID) ([]model.BlogSummary, error) {
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
				SELECT bp.photo_url FROM
				blog_photos bp WHERE bp.blog_id = b.id
				ORDER BY bp.id ASC
			),'') AS blog_thumbnail,
			COUNT(l.id) FILTER (WHERE like_type = 'like') AS likes_count,
			COUNT(l.id) FILTER (WHERE like_type = 'dislike') AS dislikes_count,
			COUNT(DISTINCT c.id) AS comments_count

		FROM bookmarks bm
		JOIN blogs b ON bm.blog_id = b.id
		LEFT JOIN likes l ON l.blog_id = b.id
		LEFT JOIN comments c ON c.blog_id = b.id
		WHERE bm.user_id = $1
		GROUP BY b.id, bm.id
		ORDER BY bm.id DESC
	`

	var blogs []model.BlogSummary

	rows, err := u.DB.Query(query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var blog model.BlogSummary

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

	return blogs, nil
}
