// backend/repositories/user_repository.go
package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shimaf4979/pamfree-backend/models"
)

// UserRepository はユーザーデータへのアクセスを提供するインターフェース
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*models.User, error)
}

// MySQLUserRepository はMySQLデータベースを使用したUserRepositoryの実装
type MySQLUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository は新しいMySQLUserRepositoryを作成する
func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &MySQLUserRepository{db: db}
}

// Create は新しいユーザーを作成する
func (r *MySQLUserRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	user.CreatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password, name, role, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
		user.CreatedAt,
	)

	return err
}

// GetByID はIDによりユーザーを取得する
func (r *MySQLUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail はメールアドレスによりユーザーを取得する
func (r *MySQLUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at
		FROM users
		WHERE email = ?
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update はユーザー情報を更新する
func (r *MySQLUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, name = ?, role = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.Name,
		user.Role,
		user.ID,
	)

	return err
}

// Delete はユーザーを削除する
func (r *MySQLUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetAll は全てのユーザーを取得する
func (r *MySQLUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
