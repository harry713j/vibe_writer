package repo

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// create user
func (r *UserRepository) CreateUser(username, email, password string) (*model.User, error) {
	user := &model.User{
		Id:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.DB.Exec("INSERT INTO users(id, username, email, password_hash, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6)",
		user.Id, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// get user by id and username
func (r *UserRepository) GetUserById(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELCT * FROM users WHERE id=$1", id).Scan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByIdentifier(idntifier string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELECT * FROM users WHERE username=$1 OR email=$2", idntifier, idntifier).Scan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELCT * FROM users WHERE username=$1", username).Scan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELCT * FROM users WHERE email=$1", email).Scan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// delete user
func (r *UserRepository) DeleteUser(userId uuid.UUID) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE user_id=$1", userId)

	return err
}
