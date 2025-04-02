// backend/repositories/pin_repository.go
package repositories

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shimaf4979/pamfree-backend/models"
)

// PinRepository はピンデータへのアクセスを提供するインターフェース
type PinRepository interface {
	Create(ctx context.Context, pin *models.Pin) error
	GetByID(ctx context.Context, id string) (*models.Pin, error)
	GetByFloorID(ctx context.Context, floorID string) ([]*models.Pin, error)
	GetByFloorIDs(ctx context.Context, floorIDs []string) ([]*models.Pin, error)
	Update(ctx context.Context, pin *models.Pin) error
	Delete(ctx context.Context, id string) error
}

// MySQLPinRepository はMySQLデータベースを使用したPinRepositoryの実装
type MySQLPinRepository struct {
	db *sql.DB
}

// NewMySQLPinRepository は新しいMySQLPinRepositoryを作成する
func NewMySQLPinRepository(db *sql.DB) PinRepository {
	return &MySQLPinRepository{db: db}
}

// Create は新しいピンを作成する
func (r *MySQLPinRepository) Create(ctx context.Context, pin *models.Pin) error {
	if pin.ID == "" {
		pin.ID = uuid.New().String()
	}
	pin.CreatedAt = time.Now()
	pin.UpdatedAt = time.Now()

	query := `
		INSERT INTO pins (id, floor_id, title, description, x_position, y_position, image_url, editor_id, editor_nickname, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		pin.ID,
		pin.FloorID,
		pin.Title,
		pin.Description,
		pin.XPosition,
		pin.YPosition,
		pin.ImageURL,
		pin.EditorID,
		pin.EditorNickname,
		pin.CreatedAt,
		pin.UpdatedAt,
	)

	return err
}

// GetByID はIDによりピンを取得する
func (r *MySQLPinRepository) GetByID(ctx context.Context, id string) (*models.Pin, error) {
	query := `
		SELECT id, floor_id, title, description, x_position, y_position, image_url, editor_id, editor_nickname, created_at, updated_at
		FROM pins
		WHERE id = ?
	`

	var pin models.Pin
	var editorID, editorNickname, imageURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pin.ID,
		&pin.FloorID,
		&pin.Title,
		&pin.Description,
		&pin.XPosition,
		&pin.YPosition,
		&imageURL,
		&editorID,
		&editorNickname,
		&pin.CreatedAt,
		&pin.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// NULL値の処理
	if imageURL.Valid {
		pin.ImageURL = imageURL.String
	}
	if editorID.Valid {
		pin.EditorID = editorID.String
	}
	if editorNickname.Valid {
		pin.EditorNickname = editorNickname.String
	}

	return &pin, nil
}

// GetByFloorID はフロアIDによりピンを取得する
func (r *MySQLPinRepository) GetByFloorID(ctx context.Context, floorID string) ([]*models.Pin, error) {
	query := `
		SELECT id, floor_id, title, description, x_position, y_position, image_url, editor_id, editor_nickname, created_at, updated_at
		FROM pins
		WHERE floor_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, floorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []*models.Pin
	for rows.Next() {
		var pin models.Pin
		var editorID, editorNickname, imageURL sql.NullString

		if err := rows.Scan(
			&pin.ID,
			&pin.FloorID,
			&pin.Title,
			&pin.Description,
			&pin.XPosition,
			&pin.YPosition,
			&imageURL,
			&editorID,
			&editorNickname,
			&pin.CreatedAt,
			&pin.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// NULL値の処理
		if imageURL.Valid {
			pin.ImageURL = imageURL.String
		}
		if editorID.Valid {
			pin.EditorID = editorID.String
		}
		if editorNickname.Valid {
			pin.EditorNickname = editorNickname.String
		}

		pins = append(pins, &pin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pins, nil
}

// GetByFloorIDs は複数のフロアIDに対応するピンを取得する
func (r *MySQLPinRepository) GetByFloorIDs(ctx context.Context, floorIDs []string) ([]*models.Pin, error) {
	// フロアIDが空の場合は空のスライスを返す
	if len(floorIDs) == 0 {
		return []*models.Pin{}, nil
	}

	// プレースホルダーを作成 (IN句用)
	placeholders := make([]string, len(floorIDs))
	args := make([]interface{}, len(floorIDs))
	for i, id := range floorIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	// SQLクエリを構築
	query := `
		SELECT id, floor_id, title, description, x_position, y_position, image_url, editor_id, editor_nickname, created_at, updated_at
		FROM pins
		WHERE floor_id IN (` + strings.Join(placeholders, ",") + `)
		ORDER BY created_at ASC
	`

	// クエリを実行
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 結果を処理
	var pins []*models.Pin
	for rows.Next() {
		var pin models.Pin
		var editorID, editorNickname, imageURL sql.NullString

		if err := rows.Scan(
			&pin.ID,
			&pin.FloorID,
			&pin.Title,
			&pin.Description,
			&pin.XPosition,
			&pin.YPosition,
			&imageURL,
			&editorID,
			&editorNickname,
			&pin.CreatedAt,
			&pin.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// NULLの場合は空文字に
		if imageURL.Valid {
			pin.ImageURL = imageURL.String
		}
		if editorID.Valid {
			pin.EditorID = editorID.String
		}
		if editorNickname.Valid {
			pin.EditorNickname = editorNickname.String
		}

		pins = append(pins, &pin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pins, nil
}

// Update はピン情報を更新する
func (r *MySQLPinRepository) Update(ctx context.Context, pin *models.Pin) error {
	pin.UpdatedAt = time.Now()

	query := `
		UPDATE pins
		SET title = ?, description = ?, x_position = ?, y_position = ?, image_url = ?, 
		    editor_id = ?, editor_nickname = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		pin.Title,
		pin.Description,
		pin.XPosition,
		pin.YPosition,
		pin.ImageURL,
		pin.EditorID,
		pin.EditorNickname,
		pin.UpdatedAt,
		pin.ID,
	)

	return err
}

// Delete はピンを削除する
func (r *MySQLPinRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM pins WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
