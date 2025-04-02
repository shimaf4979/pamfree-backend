// backend/repositories/public_editor_repository.go
package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shimaf4979/pamfree-backend/models"
)

// PublicEditorRepository は公開編集者データへのアクセスを提供するインターフェース
type PublicEditorRepository interface {
	Create(ctx context.Context, editor *models.PublicEditor) error
	GetByID(ctx context.Context, id string) (*models.PublicEditor, error)
	GetByToken(ctx context.Context, token string) (*models.PublicEditor, error)
	GetByMapID(ctx context.Context, mapID string) ([]*models.PublicEditor, error)
	Update(ctx context.Context, editor *models.PublicEditor) error
	UpdateLastActive(ctx context.Context, id string) error
}

// MySQLPublicEditorRepository はMySQLデータベースを使用したPublicEditorRepositoryの実装
type MySQLPublicEditorRepository struct {
	db *sql.DB
}

// NewMySQLPublicEditorRepository は新しいMySQLPublicEditorRepositoryを作成する
func NewMySQLPublicEditorRepository(db *sql.DB) PublicEditorRepository {
	return &MySQLPublicEditorRepository{db: db}
}

// Create は新しい公開編集者を作成する
func (r *MySQLPublicEditorRepository) Create(ctx context.Context, editor *models.PublicEditor) error {
	if editor.ID == "" {
		editor.ID = uuid.New().String()
	}
	editor.CreatedAt = time.Now()
	editor.LastActive = time.Now()

	query := `
		INSERT INTO public_editors (id, map_id, nickname, editor_token, created_at, last_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		editor.ID,
		editor.MapID,
		editor.Nickname,
		editor.EditorToken,
		editor.CreatedAt,
		editor.LastActive,
	)

	return err
}

// GetByID はIDにより公開編集者を取得する
func (r *MySQLPublicEditorRepository) GetByID(ctx context.Context, id string) (*models.PublicEditor, error) {
	query := `
		SELECT id, map_id, nickname, editor_token, created_at, last_active
		FROM public_editors
		WHERE id = ?
	`

	var editor models.PublicEditor
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&editor.ID,
		&editor.MapID,
		&editor.Nickname,
		&editor.EditorToken,
		&editor.CreatedAt,
		&editor.LastActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &editor, nil
}

// GetByToken はトークンにより公開編集者を取得する
func (r *MySQLPublicEditorRepository) GetByToken(ctx context.Context, token string) (*models.PublicEditor, error) {
	query := `
		SELECT id, map_id, nickname, editor_token, created_at, last_active
		FROM public_editors
		WHERE editor_token = ?
	`

	var editor models.PublicEditor
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&editor.ID,
		&editor.MapID,
		&editor.Nickname,
		&editor.EditorToken,
		&editor.CreatedAt,
		&editor.LastActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &editor, nil
}

// GetByMapID はマップIDにより公開編集者を取得する
func (r *MySQLPublicEditorRepository) GetByMapID(ctx context.Context, mapID string) ([]*models.PublicEditor, error) {
	query := `
		SELECT id, map_id, nickname, editor_token, created_at, last_active
		FROM public_editors
		WHERE map_id = ?
		ORDER BY last_active DESC
	`

	rows, err := r.db.QueryContext(ctx, query, mapID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var editors []*models.PublicEditor
	for rows.Next() {
		var editor models.PublicEditor
		if err := rows.Scan(
			&editor.ID,
			&editor.MapID,
			&editor.Nickname,
			&editor.EditorToken,
			&editor.CreatedAt,
			&editor.LastActive,
		); err != nil {
			return nil, err
		}
		editors = append(editors, &editor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return editors, nil
}

// Update は公開編集者情報を更新する
func (r *MySQLPublicEditorRepository) Update(ctx context.Context, editor *models.PublicEditor) error {
	query := `
		UPDATE public_editors
		SET nickname = ?, editor_token = ?, last_active = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		editor.Nickname,
		editor.EditorToken,
		editor.LastActive,
		editor.ID,
	)

	return err
}

// UpdateLastActive は最終アクティブ時間を更新する
func (r *MySQLPublicEditorRepository) UpdateLastActive(ctx context.Context, id string) error {
	query := `
		UPDATE public_editors
		SET last_active = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		time.Now(),
		id,
	)

	return err
}
