// backend/services/viewer_service.go
package services

import (
	"context"
	"errors"

	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/repositories"
)

// ViewerService はビューワー機能に関する操作を提供するインターフェース
type ViewerService interface {
	GetMapData(ctx context.Context, mapID string) (*models.ViewerData, error)
}

// DefaultViewerService はViewerServiceの実装
type DefaultViewerService struct {
	mapRepo   repositories.MapRepository
	floorRepo repositories.FloorRepository
	pinRepo   repositories.PinRepository
}

// NewViewerService は新しいViewerServiceを作成する
func NewViewerService(
	mapRepo repositories.MapRepository,
	floorRepo repositories.FloorRepository,
	pinRepo repositories.PinRepository,
) ViewerService {
	return &DefaultViewerService{
		mapRepo:   mapRepo,
		floorRepo: floorRepo,
		pinRepo:   pinRepo,
	}
}

// GetMapData はマップの全データを取得する
func (s *DefaultViewerService) GetMapData(ctx context.Context, mapID string) (*models.ViewerData, error) {
	// マップデータを取得
	mapData, err := s.mapRepo.GetByID(ctx, mapID)
	if err != nil {
		return nil, err
	}
	if mapData == nil {
		return nil, errors.New("マップが見つかりません")
	}

	// フロアデータを取得
	floors, err := s.floorRepo.GetByMapID(ctx, mapData.ID)
	if err != nil {
		return nil, err
	}

	// ピンデータを取得
	var pins []*models.Pin
	if len(floors) > 0 {
		// フロアIDのスライスを作成
		floorIDs := make([]string, len(floors))
		for i, floor := range floors {
			floorIDs[i] = floor.ID
		}

		// 全フロアのピンを取得
		pins, err = s.pinRepo.GetByFloorIDs(ctx, floorIDs)
		if err != nil {
			return nil, err
		}
	}

	// ビューワーデータを構築
	viewerData := &models.ViewerData{
		Map:    mapData,
		Floors: floors,
		Pins:   pins,
	}

	return viewerData, nil
}
