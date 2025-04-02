// backend/services/floor_service.go
package services

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/yourname/mapapp/models"
	"github.com/yourname/mapapp/repositories"
)

// FloorService はフロア関連のビジネスロジックを提供するインターフェース
type FloorService interface {
	CreateFloor(ctx context.Context, req models.FloorCreate, userID string) (*models.Floor, error)
	GetFloorsByMapID(ctx context.Context, mapID string) ([]*models.Floor, error)
	GetFloorByID(ctx context.Context, id string) (*models.Floor, error)
	UpdateFloor(ctx context.Context, id string, req models.FloorUpdate, userID string) (*models.Floor, error)
	DeleteFloor(ctx context.Context, id string, userID string) error
	UploadFloorImage(ctx context.Context, floorID string, userID string, file io.Reader) (string, error)
}

// FloorServiceImpl はFloorServiceの実装
type FloorServiceImpl struct {
	floorRepo  repositories.FloorRepository
	mapRepo    repositories.MapRepository
	cloudinary *cloudinary.Cloudinary
}

// NewFloorService は新しいFloorServiceを作成する
func NewFloorService(
	floorRepo repositories.FloorRepository,
	mapRepo repositories.MapRepository,
	cloudinary *cloudinary.Cloudinary,
) FloorService {
	return &FloorServiceImpl{
		floorRepo:  floorRepo,
		mapRepo:    mapRepo,
		cloudinary: cloudinary,
	}
}

// CreateFloor は新しいフロアを作成する
func (s *FloorServiceImpl) CreateFloor(ctx context.Context, req models.FloorCreate, userID string) (*models.Floor, error) {
	// マップを取得して所有者を確認
	mapObj, err := s.mapRepo.GetByID(ctx, req.MapID)
	if err != nil {
		return nil, err
	}

	if mapObj == nil {
		return nil, errors.New("マップが見つかりません")
	}

	if mapObj.UserID != userID {
		return nil, errors.New("このマップを編集する権限がありません")
	}

	// 新しいフロアを作成
	floor := &models.Floor{
		MapID:       req.MapID,
		FloorNumber: req.FloorNumber,
		Name:        req.Name,
	}

	if err := s.floorRepo.Create(ctx, floor); err != nil {
		return nil, err
	}

	return floor, nil
}

// GetFloorsByMapID はマップに属するすべてのフロアを取得する
func (s *FloorServiceImpl) GetFloorsByMapID(ctx context.Context, mapID string) ([]*models.Floor, error) {
	// マップIDを使用してフロアを取得
	return s.floorRepo.GetByMapID(ctx, mapID)
}

// GetFloorByID はIDによりフロアを取得する
func (s *FloorServiceImpl) GetFloorByID(ctx context.Context, id string) (*models.Floor, error) {
	return s.floorRepo.GetByID(ctx, id)
}

// UpdateFloor はフロア情報を更新する
func (s *FloorServiceImpl) UpdateFloor(ctx context.Context, id string, req models.FloorUpdate, userID string) (*models.Floor, error) {
	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if floor == nil {
		return nil, errors.New("フロアが見つかりません")
	}

	// マップを取得して所有者を確認
	mapObj, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return nil, err
	}

	if mapObj == nil {
		return nil, errors.New("マップが見つかりません")
	}

	if mapObj.UserID != userID {
		return nil, errors.New("このフロアを編集する権限がありません")
	}

	// フロア情報を更新
	if req.Name != "" {
		floor.Name = req.Name
	}

	if req.ImageURL != "" {
		floor.ImageURL = req.ImageURL
	}

	if err := s.floorRepo.Update(ctx, floor); err != nil {
		return nil, err
	}

	return floor, nil
}

// DeleteFloor はフロアを削除する
func (s *FloorServiceImpl) DeleteFloor(ctx context.Context, id string, userID string) error {
	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if floor == nil {
		return errors.New("フロアが見つかりません")
	}

	// マップを取得して所有者を確認
	mapObj, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return err
	}

	if mapObj == nil {
		return errors.New("マップが見つかりません")
	}

	if mapObj.UserID != userID {
		return errors.New("このフロアを削除する権限がありません")
	}

	// フロアを削除
	return s.floorRepo.Delete(ctx, id)
}

// UploadFloorImage はフロアの画像をアップロードする
func (s *FloorServiceImpl) UploadFloorImage(ctx context.Context, floorID string, userID string, file io.Reader) (string, error) {
	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, floorID)
	if err != nil {
		return "", err
	}

	if floor == nil {
		return "", errors.New("フロアが見つかりません")
	}

	// マップを取得して所有者を確認
	mapObj, err := s.mapRepo.GetByID(ctx, floor.MapID)
	if err != nil {
		return "", err
	}

	if mapObj == nil {
		return "", errors.New("マップが見つかりません")
	}

	if mapObj.UserID != userID {
		return "", errors.New("このフロアを編集する権限がありません")
	}

	// Cloudinaryにアップロード
	uploadResult, err := s.cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "map_images",
		Format:         "webp",
		Transformation: "q_auto",
		ResourceType:   "image",
		PublicID:       fmt.Sprintf("%s_%s", mapObj.MapID, floorID),
	})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}