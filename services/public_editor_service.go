// backend/services/public_editor_service.go
package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/repositories"
)

// PublicEditorService は公開編集者に関する操作を提供するインターフェース
type PublicEditorService interface {
	Register(ctx context.Context, mapID, nickname string) (*models.PublicEditor, error)
	Verify(ctx context.Context, editorID, token string) (*models.PublicEditor, error)
	GetByID(ctx context.Context, editorID string) (*models.PublicEditor, error)
	GetByMapID(ctx context.Context, mapID string) ([]*models.PublicEditor, error)
	UpdateLastActive(ctx context.Context, editorID string) error
}

// DefaultPublicEditorService はPublicEditorServiceの実装
type DefaultPublicEditorService struct {
	publicEditorRepo repositories.PublicEditorRepository
	mapRepo          repositories.MapRepository
}

// NewPublicEditorService は新しいPublicEditorServiceを作成する
func NewPublicEditorService(
	publicEditorRepo repositories.PublicEditorRepository,
	mapRepo repositories.MapRepository,
) PublicEditorService {
	return &DefaultPublicEditorService{
		publicEditorRepo: publicEditorRepo,
		mapRepo:          mapRepo,
	}
}

// Register は新しい公開編集者を登録する
func (s *DefaultPublicEditorService) Register(ctx context.Context, mapID, nickname string) (*models.PublicEditor, error) {
	// マップの存在と公開編集可能性を確認
	mapData, err := s.mapRepo.GetByID(ctx, mapID)
	if err != nil {
		return nil, err
	}
	if mapData == nil {
		return nil, errors.New("マップが見つかりません")
	}
	if !mapData.IsPubliclyEditable {
		return nil, errors.New("このマップは公開編集が許可されていません")
	}

	// ランダムなトークンを生成
	token, err := generateToken(32)
	if err != nil {
		return nil, err
	}

	// 公開編集者を作成
	editor := &models.PublicEditor{
		ID:          uuid.New().String(),
		MapID:       mapID,
		Nickname:    nickname,
		EditorToken: token,
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
	}

	// リポジトリに保存
	if err := s.publicEditorRepo.Create(ctx, editor); err != nil {
		return nil, err
	}

	return editor, nil
}

// Verify は編集者トークンを検証する
func (s *DefaultPublicEditorService) Verify(ctx context.Context, editorID, token string) (*models.PublicEditor, error) {
	editor, err := s.publicEditorRepo.GetByID(ctx, editorID)
	if err != nil {
		return nil, err
	}
	if editor == nil {
		return nil, errors.New("編集者が見つかりません")
	}

	// トークンの検証
	if editor.EditorToken != token {
		return nil, errors.New("無効なトークンです")
	}

	// 最終アクティブ時間を更新
	if err := s.UpdateLastActive(ctx, editorID); err != nil {
		return nil, err
	}

	return editor, nil
}

// GetByID はIDによって公開編集者を取得する
func (s *DefaultPublicEditorService) GetByID(ctx context.Context, editorID string) (*models.PublicEditor, error) {
	return s.publicEditorRepo.GetByID(ctx, editorID)
}

// GetByMapID はマップIDによって公開編集者を取得する
func (s *DefaultPublicEditorService) GetByMapID(ctx context.Context, mapID string) ([]*models.PublicEditor, error) {
	return s.publicEditorRepo.GetByMapID(ctx, mapID)
}

// UpdateLastActive は公開編集者の最終アクティブ時間を更新する
func (s *DefaultPublicEditorService) UpdateLastActive(ctx context.Context, editorID string) error {
	editor, err := s.publicEditorRepo.GetByID(ctx, editorID)
	if err != nil {
		return err
	}
	if editor == nil {
		return errors.New("編集者が見つかりません")
	}

	editor.LastActive = time.Now()
	return s.publicEditorRepo.Update(ctx, editor)
}

// generateToken はランダムなトークンを生成する
func generateToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
