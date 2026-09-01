package repository

import (
	"context"
	"fmt"

	"go101/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, name, email string) (*model.User, error) {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at`

	var user model.User
	if err := r.db.QueryRow(ctx, query, name, email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1`

	var user model.User
	if err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM users
		ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return users, nil
}
