package repository

import (
	"context"
	"database/sql"

	"chow/internal/model"

	"github.com/google/uuid"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, data *model.User) (*model.User, error) {
	var user model.User
	query := `
        INSERT INTO users(email, username, password, role)
        VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, password, role, created_at, updated_at
    `
	if err := r.db.QueryRowContext(ctx, query, data.Email, data.Username, data.Password, data.Role).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmailOrUsername retrieves a user by their email or username. The fallback is email if key does not match `username`
func (r *UserRepository) GetUserByEmailOrUsername(ctx context.Context, key, value string) (*model.User, error) {
	query := `
		SELECT id, email, username, password, role, created_at, updated_at
		FROM users 
		WHERE email = $1
		`

	if key == "username" {
		query = `
			SELECT id, email, username, password, role, created_at, updated_at
			FROM users 
			WHERE username = $1
		`
	}

	var user model.User
	err := r.db.QueryRow(query, value).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, username, password, role, created_at, updated_at
		FROM users  
		WHERE id = $1
		`
	var user model.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateByID(ctx context.Context, id uuid.UUID, data *model.User) (*model.User, error) {
	var user model.User
	query := `
		UPDATE users
		SET email = $1, username = $2, password = $3, role = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, email, username, password, role, created_at, updated_at
    `
	if err := r.db.QueryRowContext(ctx, query, data.Email, data.Username, data.Password, data.Role).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
