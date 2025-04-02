// backend/services/public_editor_service.go
package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/yourname/mapapp/models"
	"github.com/yourname/mapapp/repositories"
)

// PublicEditorService は公開編集者関連のビジネスロジックを提供するインターフェース
type PublicEditorService interface {
	RegisterEditor(ctx context.Context, req models.PublicEditorRegister) (*models.PublicEditorResponse, error)
	VerifyEditorToken(ctx context.Context, req models.PublicEditorVerify) (*models.PublicEditor, error)
}

// PublicEditorServiceImpl はPublicEditorServiceの実装
type PublicEditorServiceImpl struct {
	publicEditorRepo repositories.PublicEditorRepository
	mapRepo          repositories.MapRepository
}

// NewPublicEditorService は新しいPublicEditorServiceを作成する
func NewPublicEditorService(
	publicEditorRepo repositories.PublicEditorRepository,
	mapRepo repositories.MapRepository,
) PublicEditorService {
	return &PublicEditorServiceImpl{
		publicEditorRepo: publicEditorRepo,
		mapRepo:          mapRepo,
	}
}

// RegisterEditor は新しい公開編集者を登録する
func (s *PublicEditorServiceImpl) RegisterEditor(ctx context.Context, req models.PublicEditorRegister) (*models.PublicEditorResponse, error) {
	// マップの存在と公開編集可能かを確認
	mapObj, err := s.mapRepo.GetByMapID(ctx, req.MapID)
	if err != nil {
		return nil, err
	}

	if mapObj == nil {
		return nil, errors.New("マップが見つかりません")
	}

	if !mapObj.IsPubliclyEditable {
		return nil, errors.New("このマップは公開編集が有効になっていません")
	}

	// ランダムなトークンを生成
	token, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	// 公開編集者を作成
	editor := &models.PublicEditor{
		MapID:       mapObj.ID,
		Nickname:    req.Nickname,
		EditorToken: token,
	}

	if err := s.publicEditorRepo.Create(ctx, editor); err != nil {
		return nil, err
	}

	// レスポンスを作成
	response := &models.PublicEditorResponse{
		ID:       editor.ID,
		Nickname: editor.Nickname,
		MapID:    mapObj.MapID,
		Token:    editor.EditorToken,
	}

	return response, nil
}

// VerifyEditorToken は公開編集者トークンを検証する
func (s *PublicEditorServiceImpl) VerifyEditorToken(ctx context.Context, req models.PublicEditorVerify) (*models.PublicEditor, error) {
	// 公開編集者を取得
	editor, err := s.publicEditorRepo.GetByID(ctx, req.EditorID)
	if err != nil {
		return nil, err
	}

	if editor == nil {
		return nil, errors.New("編集者が見つかりません")
	}

	// トークンを検証
	if editor.EditorToken != req.Token {
		return nil, errors.New("無効なトークンです")
	}

	// 最終アクティブ時間を更新
	if err := s.publicEditorRepo.UpdateLastActive(ctx, editor.ID); err != nil {
		return nil, err
	}

	return editor, nil
}

// generateRandomToken はランダムなトークンを生成する
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}