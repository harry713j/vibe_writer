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

	query := `SELECT up.full_name, up.bio, up.avatar_url, u.id, u.username, u.email, u.created_at, u.updated_at 
	FROM user_profiles up INNER JOIN users u ON up.user_id = u.id 
	WHERE up.user_id = $1`

	err := u.DB.QueryRow(query, userId).Scan(&userData)

	if err != nil {
		return nil, err
	}

	return &userData, nil
}
