// services/map_service.go
package services

import (
	"context"
	"errors"

	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/repositories"
)

// MapService マップに関する操作を提供するインターフェース
type MapService interface {
	GetMapsByUserID(ctx context.Context, userID string) ([]*models.Map, error)
	GetMapByID(ctx context.Context, id string) (*models.Map, error)
	CreateMap(ctx context.Context, m *models.Map) error
	UpdateMap(ctx context.Context, m *models.Map) error
	DeleteMap(ctx context.Context, id string) error
}

// DefaultMapService はMapServiceの実装
type DefaultMapService struct {
	mapRepo repositories.MapRepository
}

// NewMapService 新しいMapServiceを作成
func NewMapService(mapRepo repositories.MapRepository) MapService {
	return &DefaultMapService{
		mapRepo: mapRepo,
	}
}

// GetMapsByUserID ユーザーIDによるマップ一覧の取得
func (s *DefaultMapService) GetMapsByUserID(ctx context.Context, userID string) ([]*models.Map, error) {
	return s.mapRepo.GetByUserID(ctx, userID)
}

// GetMapByID IDによるマップの取得
func (s *DefaultMapService) GetMapByID(ctx context.Context, id string) (*models.Map, error) {
	return s.mapRepo.GetByID(ctx, id)
}

// CreateMap マップの作成
func (s *DefaultMapService) CreateMap(ctx context.Context, m *models.Map) error {
	// IDの重複チェック
	existingMap, err := s.mapRepo.GetByID(ctx, m.ID)
	if err != nil {
		return err
	}
	if existingMap != nil {
		return errors.New("このIDは既に使用されています")
	}

	return s.mapRepo.Create(ctx, m)
}

// UpdateMap マップの更新
func (s *DefaultMapService) UpdateMap(ctx context.Context, m *models.Map) error {
	// マップが存在するか確認
	existingMap, err := s.mapRepo.GetByID(ctx, m.ID)
	if err != nil {
		return err
	}
	if existingMap == nil {
		return errors.New("マップが見つかりません")
	}

	return s.mapRepo.Update(ctx, m)
}

// DeleteMap マップの削除
func (s *DefaultMapService) DeleteMap(ctx context.Context, id string) error {
	// マップが存在するか確認
	existingMap, err := s.mapRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingMap == nil {
		return errors.New("マップが見つかりません")
	}

	return s.mapRepo.Delete(ctx, id)
}
