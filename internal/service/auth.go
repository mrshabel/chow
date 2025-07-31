package service

import (
	"chow/internal/config"
	"chow/internal/model"
	"chow/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// errors
var (
	ErrAlreadyExist       = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	cfg      *config.Config
	userRepo *repository.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

// Register
func (s *AuthService) Register(ctx context.Context, data *model.User) (*model.User, error) {
	// check if user exists
	_, err := s.userRepo.GetUserByEmailOrUsername(ctx, "email", data.Email)
	if err == nil {
		return nil, ErrAlreadyExist
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	// hash password
	data.Password, err = s.hashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	// set default user role
	data.Role = model.AppUser

	// create user
	return s.userRepo.Create(ctx, data)
}

// Login authenticates a user and generate their access token
func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	// check if user exists
	user, err := s.userRepo.GetUserByEmailOrUsername(ctx, "email", email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// verify password
	if err := s.verifyPassword(password, user.Password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// generate access token
	token, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, err
}

// forgot password
// func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error)

// reset password

// token helpers (generate and validate)

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	now := time.Now()
	expiry := now.Add(s.cfg.JWTExpiryMinutes)

	// create token with claims
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"iat":      now.Unix(),
		"exp":      expiry.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token with secret key
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *AuthService) ValidateToken(token string) (jwt.MapClaims, error) {
	// parse token
	parsedToken, err := jwt.Parse(token, func(parsedToken *jwt.Token) (any, error) {
		// validate signing method
		if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.cfg.JWTSecret, nil
	})
	if err != nil {
		// check for expiry
		if err == jwt.ErrTokenExpired {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// extract claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// password utils

func (s *AuthService) hashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), err
}

func (s *AuthService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
