// repositories/map_repository.go
package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/shimaf4979/pamfree-backend/models"
)

// MapRepository マップデータへのアクセスを提供するインターフェース
type MapRepository interface {
	Create(ctx context.Context, m *models.Map) error
	GetByID(ctx context.Context, id string) (*models.Map, error)
	GetByMapID(ctx context.Context, mapID string) (*models.Map, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Map, error)
	Update(ctx context.Context, m *models.Map) error
	Delete(ctx context.Context, id string) error
}

// MySQLMapRepository はMySQLデータベースを使用したMapRepositoryの実装
type MySQLMapRepository struct {
	db *sql.DB
}

// NewMySQLMapRepository は新しいMySQLMapRepositoryを作成する
func NewMySQLMapRepository(db *sql.DB) MapRepository {
	return &MySQLMapRepository{db: db}
}

// Create は新しいマップを作成する
func (r *MySQLMapRepository) Create(ctx context.Context, m *models.Map) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	query := `
		INSERT INTO maps (id, map_id, title, description, user_id, is_publicly_editable, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		m.ID,
		m.MapID,
		m.Title,
		m.Description,
		m.UserID,
		m.IsPubliclyEditable,
		m.CreatedAt,
		m.UpdatedAt,
	)

	return err
}

// GetByID はIDによりマップを取得する
func (r *MySQLMapRepository) GetByID(ctx context.Context, id string) (*models.Map, error) {
	query := `
		SELECT id, map_id, title, description, user_id, is_publicly_editable, created_at, updated_at
		FROM maps
		WHERE id = ?
	`

	var m models.Map
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.MapID,
		&m.Title,
		&m.Description,
		&m.UserID,
		&m.IsPubliclyEditable,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &m, nil
}

// GetByMapID はマップIDによりマップを取得する
func (r *MySQLMapRepository) GetByMapID(ctx context.Context, mapID string) (*models.Map, error) {
	query := `
		SELECT id, map_id, title, description, user_id, is_publicly_editable, created_at, updated_at
		FROM maps
		WHERE map_id = ?
	`

	var m models.Map
	err := r.db.QueryRowContext(ctx, query, mapID).Scan(
		&m.ID,
		&m.MapID,
		&m.Title,
		&m.Description,
		&m.UserID,
		&m.IsPubliclyEditable,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &m, nil
}

// GetByUserID はユーザーIDによりマップ一覧を取得する
func (r *MySQLMapRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Map, error) {
	query := `
		SELECT id, map_id, title, description, user_id, is_publicly_editable, created_at, updated_at
		FROM maps
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var maps []*models.Map
	for rows.Next() {
		var m models.Map
		if err := rows.Scan(
			&m.ID,
			&m.MapID,
			&m.Title,
			&m.Description,
			&m.UserID,
			&m.IsPubliclyEditable,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		maps = append(maps, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return maps, nil
}

// Update はマップ情報を更新する
func (r *MySQLMapRepository) Update(ctx context.Context, m *models.Map) error {
	m.UpdatedAt = time.Now()

	query := `
		UPDATE maps
		SET title = ?, description = ?, is_publicly_editable = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		m.Title,
		m.Description,
		m.IsPubliclyEditable,
		m.UpdatedAt,
		m.ID,
	)

	return err
}

// Delete はマップを削除する
func (r *MySQLMapRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM maps WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
