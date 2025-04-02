// services/map_service.go
package services

import (
	"context"

	"github.com/yourname/mapapp/models"
	"github.com/yourname/mapapp/repositories"
)

// MapService マップサービス
type MapService struct {
	mapRepo repositories.MapRepository
}

// NewMapService 新しいマップサービスを作成
func NewMapService(mapRepo repositories.MapRepository) *MapService {
	return &MapService{
		mapRepo: mapRepo,
	}
}

// GetMapsByUserID ユーザーIDによるマップ一覧の取得
func (s *MapService) GetMapsByUserID(ctx context.Context, userID string) ([]*models.Map, error) {
	return s.mapRepo.GetByUserID(ctx, userID)
}

// GetMapByID IDによるマップの取得
func (s *MapService) GetMapByID(ctx context.Context, id string) (*models.Map, error) {
	return s.mapRepo.GetByID(ctx, id)
}

// GetMapByMapID マップIDによるマップの取得
func (s *MapService) GetMapByMapID(ctx context.Context, mapID string) (*models.Map, error) {
	return s.mapRepo.GetByMapID(ctx, mapID)
}

// MapIDExists マップIDが既に存在するか確認
func (s *MapService) MapIDExists(ctx context.Context, mapID string) (bool, error) {
	m, err := s.mapRepo.GetByMapID(ctx, mapID)
	if err != nil {
		return false, err
	}
	return m != nil, nil
}

// CreateMap マップの作成
func (s *MapService) CreateMap(ctx context.Context, m *models.Map) error {
	return s.mapRepo.Create(ctx, m)
}

// UpdateMap マップの更新
func (s *MapService) UpdateMap(ctx context.Context, m *models.Map) error {
	return s.mapRepo.Update(ctx, m)
}

// DeleteMap マップの削除
func (s *MapService) DeleteMap(ctx context.Context, id string) error {
	return s.mapRepo.Delete(ctx, id)
}