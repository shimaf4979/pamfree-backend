// backend/services/pin_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/repositories"
)

// PinService はピンに関する操作を提供するインターフェース
type PinService interface {
	Create(ctx context.Context, userID string, input *models.PinCreate) (*models.Pin, error)
	CreatePublic(ctx context.Context, input *models.PinCreate) (*models.Pin, error)
	GetByID(ctx context.Context, id string) (*models.Pin, error)
	GetByFloorID(ctx context.Context, floorID string) ([]*models.Pin, error)
	GetByFloorIDs(ctx context.Context, floorIDs []string) ([]*models.Pin, error)
	Update(ctx context.Context, userID string, id string, input *models.PinUpdate) (*models.Pin, error)
	UpdatePublic(ctx context.Context, editorID string, id string, input *models.PinUpdate) (*models.Pin, error)
	Delete(ctx context.Context, userID string, id string) error
	DeletePublic(ctx context.Context, editorID string, id string) error
}

// DefaultPinService はPinServiceの実装
type DefaultPinService struct {
	pinRepo   repositories.PinRepository
	floorRepo repositories.FloorRepository
	mapRepo   repositories.MapRepository
}

// NewPinService は新しいPinServiceを作成する
func NewPinService(
	pinRepo repositories.PinRepository,
	floorRepo repositories.FloorRepository,
	mapRepo repositories.MapRepository,
) PinService {
	return &DefaultPinService{
		pinRepo:   pinRepo,
		floorRepo: floorRepo,
		mapRepo:   mapRepo,
	}
}

// Create は新しいピンを作成する
func (s *DefaultPinService) Create(ctx context.Context, userID string, input *models.PinCreate) (*models.Pin, error) {
	// フロアが存在するか確認
	floor, err := s.floorRepo.GetByID(ctx, input.FloorID)
	if err != nil {
		return nil, err
	}
	if floor == nil {
		return nil, errors.New("フロアが見つかりません")
	}

	// マップの所有者を確認
	map_, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return nil, err
	}
	if map_ == nil {
		return nil, errors.New("マップが見つかりません")
	}

	// 権限チェック
	if map_.UserID != userID {
		return nil, errors.New("このマップを編集する権限がありません")
	}

	// 新しいピンを作成
	pin := &models.Pin{
		ID:             uuid.New().String(),
		FloorID:        input.FloorID,
		Title:          input.Title,
		Description:    input.Description,
		XPosition:      input.XPosition,
		YPosition:      input.YPosition,
		ImageURL:       input.ImageURL,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		EditorID:       userID,
		EditorNickname: "管理者", // 管理者によるピン作成
	}

	// リポジトリに保存
	if err := s.pinRepo.Create(ctx, pin); err != nil {
		return nil, err
	}

	return pin, nil
}

// CreatePublic は公開編集用の新しいピンを作成する
func (s *DefaultPinService) CreatePublic(ctx context.Context, input *models.PinCreate) (*models.Pin, error) {
	// フロアが存在するか確認
	floor, err := s.floorRepo.GetByID(ctx, input.FloorID)
	if err != nil {
		return nil, err
	}
	if floor == nil {
		return nil, errors.New("フロアが見つかりません")
	}

	// マップが公開編集可能か確認
	map_, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return nil, err
	}
	if map_ == nil {
		return nil, errors.New("マップが見つかりません")
	}

	if !map_.IsPubliclyEditable {
		return nil, errors.New("このマップは公開編集が許可されていません")
	}

	// 編集者情報のチェック
	if input.EditorID == "" || input.EditorNickname == "" {
		return nil, errors.New("編集者情報が必要です")
	}

	// 新しいピンを作成
	pin := &models.Pin{
		ID:             uuid.New().String(),
		FloorID:        input.FloorID,
		Title:          input.Title,
		Description:    input.Description,
		XPosition:      input.XPosition,
		YPosition:      input.YPosition,
		ImageURL:       input.ImageURL,
		EditorID:       input.EditorID,
		EditorNickname: input.EditorNickname,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// リポジトリに保存
	if err := s.pinRepo.Create(ctx, pin); err != nil {
		return nil, err
	}

	return pin, nil
}

// GetByID はIDによりピンを取得する
func (s *DefaultPinService) GetByID(ctx context.Context, id string) (*models.Pin, error) {
	return s.pinRepo.GetByID(ctx, id)
}

// GetByFloorID はフロアIDによりピンを取得する
func (s *DefaultPinService) GetByFloorID(ctx context.Context, floorID string) ([]*models.Pin, error) {
	return s.pinRepo.GetByFloorID(ctx, floorID)
}

// GetByFloorIDs は複数のフロアIDに対応するピンを取得する
func (s *DefaultPinService) GetByFloorIDs(ctx context.Context, floorIDs []string) ([]*models.Pin, error) {
	return s.pinRepo.GetByFloorIDs(ctx, floorIDs)
}

// Update はピン情報を更新する
func (s *DefaultPinService) Update(ctx context.Context, userID string, id string, input *models.PinUpdate) (*models.Pin, error) {
	// ピンが存在するか確認
	pin, err := s.pinRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pin == nil {
		return nil, errors.New("ピンが見つかりません")
	}

	// フロアが存在するか確認
	floor, err := s.floorRepo.GetByID(ctx, pin.FloorID)
	if err != nil {
		return nil, err
	}
	if floor == nil {
		return nil, errors.New("フロアが見つかりません")
	}

	// マップの所有者を確認
	map_, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return nil, err
	}
	if map_ == nil {
		return nil, errors.New("マップが見つかりません")
	}

	// 権限チェック
	if map_.UserID != userID {
		return nil, errors.New("このピンを編集する権限がありません")
	}

	// ピン情報を更新
	if input.Title != "" {
		pin.Title = input.Title
	}
	if input.Description != "" {
		pin.Description = input.Description
	}
	if input.ImageURL != "" {
		pin.ImageURL = input.ImageURL
	}
	pin.UpdatedAt = time.Now()

	// リポジトリを更新
	if err := s.pinRepo.Update(ctx, pin); err != nil {
		return nil, err
	}

	return pin, nil
}

// UpdatePublic は公開編集用のピン情報を更新する
func (s *DefaultPinService) UpdatePublic(ctx context.Context, editorID string, id string, input *models.PinUpdate) (*models.Pin, error) {
	// ピンが存在するか確認
	pin, err := s.pinRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pin == nil {
		return nil, errors.New("ピンが見つかりません")
	}

	// 編集者IDを確認
	if pin.EditorID != editorID {
		return nil, errors.New("このピンを編集する権限がありません")
	}

	// フロアが存在するか確認
	floor, err := s.floorRepo.GetByID(ctx, pin.FloorID)
	if err != nil {
		return nil, err
	}
	if floor == nil {
		return nil, errors.New("フロアが見つかりません")
	}

	// マップが公開編集可能か確認
	map_, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return nil, err
	}
	if map_ == nil {
		return nil, errors.New("マップが見つかりません")
	}

	if !map_.IsPubliclyEditable {
		return nil, errors.New("このマップは公開編集が許可されていません")
	}

	// ピン情報を更新
	if input.Title != "" {
		pin.Title = input.Title
	}
	if input.Description != "" {
		pin.Description = input.Description
	}
	if input.ImageURL != "" {
		pin.ImageURL = input.ImageURL
	}
	pin.UpdatedAt = time.Now()

	// リポジトリを更新
	if err := s.pinRepo.Update(ctx, pin); err != nil {
		return nil, err
	}

	return pin, nil
}

// Delete はピンを削除する
func (s *DefaultPinService) Delete(ctx context.Context, userID string, id string) error {
	// ピンが存在するか確認
	pin, err := s.pinRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pin == nil {
		return errors.New("ピンが見つかりません")
	}

	// フロアが存在するか確認
	floor, err := s.floorRepo.GetByID(ctx, pin.FloorID)
	if err != nil {
		return err
	}
	if floor == nil {
		return errors.New("フロアが見つかりません")
	}

	// マップの所有者を確認
	map_, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return err
	}
	if map_ == nil {
		return errors.New("マップが見つかりません")
	}

	// 権限チェック
	if map_.UserID != userID {
		return errors.New("このピンを削除する権限がありません")
	}

	// ピンを削除
	return s.pinRepo.Delete(ctx, id)
}

// DeletePublic は公開編集用のピンを削除する
func (s *DefaultPinService) DeletePublic(ctx context.Context, editorID string, id string) error {
	// ピンが存在するか確認
	pin, err := s.pinRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pin == nil {
		return errors.New("ピンが見つかりません")
	}

	// 編集者IDを確認
	if pin.EditorID != editorID {
		return errors.New("このピンを削除する権限がありません")
	}

	// フロアが存在するか確認
	floor, err := s.floorRepo.GetByID(ctx, pin.FloorID)
	if err != nil {
		return err
	}
	if floor == nil {
		return errors.New("フロアが見つかりません")
	}

	// マップが公開編集可能か確認
	map_, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return err
	}
	if map_ == nil {
		return errors.New("マップが見つかりません")
	}

	if !map_.IsPubliclyEditable {
		return errors.New("このマップは公開編集が許可されていません")
	}

	// ピンを削除
	return s.pinRepo.Delete(ctx, id)
}
