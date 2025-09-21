package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
	"github.com/harry713j/vibe_writer/internal/repo"
)

type AuthService struct {
	userRepo         *repo.UserRepository
	refreshTokenRepo *repo.RefreshTokenRepository
	jwtSecret        []byte
	accessTokenTTL   time.Duration
}

func NewAuthService(userRepo *repo.UserRepository, refreshTokenRepo *repo.RefreshTokenRepository,
	jwtSecret string, accessTokenTTL time.Duration) *AuthService {

	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        []byte(jwtSecret),
		accessTokenTTL:   accessTokenTTL,
	}
}

func (service *AuthService) RegisterUser(username, email, password string) (*model.User, error) {
	// check if the user is already exists
	_, err := service.userRepo.GetUserByUsername(username)

	if err == nil {
		return nil, errors.New("Username already exists")
	}

	_, err = service.userRepo.GetUserByEmail(email)

	if err == nil {
		return nil, errors.New("Email already in use")
	}

	// create user
	hashedPassword, err := HashPassword(password)

	if err != nil {
		return nil, err
	}

	user, err := service.userRepo.CreateUser(username, email, hashedPassword)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *AuthService) LoginUser(identifier, password string) (accessToken string, refreshToken string, err error) {
	// identifier can be username or email
	user, err := service.userRepo.GetUserByIdentifier(identifier)

	if err != nil {
		return "", "", err
	}
	// create refresh token
	refresh, err := service.refreshTokenRepo.CreateRefreshToken(user.Id)

	if err != nil {
		return "", "", err
	}
	// generate the access token
	accessToken, err = service.generateAccessToken(user)

	if err != nil {
		return "", "", err
	}

	return accessToken, refresh.Token.String(), nil
}

func (service *AuthService) LogoutUser(userId uuid.UUID) error {
	err := service.refreshTokenRepo.DeleteRefreshToken(userId)

	return err
}

func (service *AuthService) RefreshAccessToken(refreshTokenStr string) (string, error) {
	// get the corresponding refresh token
	refreshTokenUUID, err := uuid.Parse(refreshTokenStr)

	if err != nil {
		return "", err
	}

	refreshToken, err := service.refreshTokenRepo.GetRefreshToken(refreshTokenUUID)

	if err != nil {
		return "", err
	}

	// check validity of the token
	if time.Now().After(refreshToken.ExpireAt) {
		return "", errors.New("Refresh token is expired")
	}
	// get the user
	user, err := service.userRepo.GetUserById(refreshToken.UserId)

	if err != nil {
		return "", err
	}
	// generate access token
	return service.generateAccessToken(user)
}

func (service *AuthService) generateAccessToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(service.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub":      user.Id.String(),
		"username": user.Username,
		"email":    user.Email,
		"iat":      time.Now().Unix(),
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(service.jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (service *AuthService) ValidateJwtToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid or Corrupted token")
		}

		return service.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("Expired Token")
		}
		return nil, errors.New("Invalid or Corrupted token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid or Corrupted token")
}
