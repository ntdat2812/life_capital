package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (email, name, password_hash, auth_provider, google_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, user.Email, user.Name, user.PasswordHash, user.AuthProvider, user.GoogleID).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, name, password_hash, auth_provider, google_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &model.User{}
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.AuthProvider, &user.GoogleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil for not found without an error
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateGoogleID(ctx context.Context, userID string, googleID string) error {
	query := `
		UPDATE users
		SET google_id = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, googleID, userID)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, name, password_hash, auth_provider, google_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &model.User{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.AuthProvider, &user.GoogleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil for not found without an error
		}
		return nil, err
	}
	return user, nil
}
