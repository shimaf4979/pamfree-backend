// backend/repositories/floor_repository.go
package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shimaf4979/pamfree-backend/models"
)

// FloorRepository はフロアデータへのアクセスを提供するインターフェース
type FloorRepository interface {
	Create(ctx context.Context, floor *models.Floor) error
	GetByID(ctx context.Context, id string) (*models.Floor, error)
	GetByMapID(ctx context.Context, mapID string) ([]*models.Floor, error)
	Update(ctx context.Context, floor *models.Floor) error
	Delete(ctx context.Context, id string) error
}

// MySQLFloorRepository はMySQLデータベースを使用したFloorRepositoryの実装
type MySQLFloorRepository struct {
	db *sql.DB
}

// NewMySQLFloorRepository は新しいMySQLFloorRepositoryを作成する
func NewMySQLFloorRepository(db *sql.DB) FloorRepository {
	return &MySQLFloorRepository{db: db}
}

// Create は新しいフロアを作成する
func (r *MySQLFloorRepository) Create(ctx context.Context, floor *models.Floor) error {
	if floor.ID == "" {
		floor.ID = uuid.New().String()
	}
	floor.CreatedAt = time.Now()
	floor.UpdatedAt = time.Now()

	query := `
		INSERT INTO floors (id, map_id, floor_number, name, image_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		floor.ID,
		floor.MapID,
		floor.FloorNumber,
		floor.Name,
		floor.ImageURL,
		floor.CreatedAt,
		floor.UpdatedAt,
	)

	return err
}

// GetByID はIDによりフロアを取得する
func (r *MySQLFloorRepository) GetByID(ctx context.Context, id string) (*models.Floor, error) {
	query := `
		SELECT id, map_id, floor_number, name, image_url, created_at, updated_at
		FROM floors
		WHERE id = ?
	`

	var floor models.Floor
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&floor.ID,
		&floor.MapID,
		&floor.FloorNumber,
		&floor.Name,
		&floor.ImageURL,
		&floor.CreatedAt,
		&floor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &floor, nil
}

// GetByMapID はマップIDによりフロアを取得する
func (r *MySQLFloorRepository) GetByMapID(ctx context.Context, mapID string) ([]*models.Floor, error) {
	query := `
		SELECT id, map_id, floor_number, name, image_url, created_at, updated_at
		FROM floors
		WHERE map_id = ?
		ORDER BY floor_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, mapID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var floors []*models.Floor
	for rows.Next() {
		var floor models.Floor
		if err := rows.Scan(
			&floor.ID,
			&floor.MapID,
			&floor.FloorNumber,
			&floor.Name,
			&floor.ImageURL,
			&floor.CreatedAt,
			&floor.UpdatedAt,
		); err != nil {
			return nil, err
		}
		floors = append(floors, &floor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return floors, nil
}

// Update はフロア情報を更新する
func (r *MySQLFloorRepository) Update(ctx context.Context, floor *models.Floor) error {
	floor.UpdatedAt = time.Now()

	query := `
		UPDATE floors
		SET name = ?, floor_number = ?, image_url = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		floor.Name,
		floor.FloorNumber,
		floor.ImageURL,
		floor.UpdatedAt,
		floor.ID,
	)

	return err
}

// Delete はフロアを削除する
func (r *MySQLFloorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM floors WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
