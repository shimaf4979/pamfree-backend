// backend/services/pin_service.go
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

// PinService はピン関連のビジネスロジックを提供するインターフェース
type PinService interface {
	CreatePin(ctx context.Context, req models.PinCreate, userID string) (*models.Pin, error)
	GetPinsByFloorID(ctx context.Context, floorID string) ([]*models.Pin, error)
	GetPinByID(ctx context.Context, id string) (*models.Pin, error)
	UpdatePin(ctx context.Context, id string, req models.PinUpdate, userID string) (*models.Pin, error)
	DeletePin(ctx context.Context, id string, userID string) error
	UploadPinImage(ctx context.Context, pinID string, userID string, file io.Reader) (string, error)
}

// PinServiceImpl はPinServiceの実装
type PinServiceImpl struct {
	pinRepo    repositories.PinRepository
	floorRepo  repositories.FloorRepository
	mapRepo    repositories.MapRepository
	cloudinary *cloudinary.Cloudinary
}

// NewPinService は新しいPinServiceを作成する
func NewPinService(
	pinRepo repositories.PinRepository,
	floorRepo repositories.FloorRepository,
	mapRepo repositories.MapRepository,
	cloudinary *cloudinary.Cloudinary,
) PinService {
	return &PinServiceImpl{
		pinRepo:    pinRepo,
		floorRepo:  floorRepo,
		mapRepo:    mapRepo,
		cloudinary: cloudinary,
	}
}

// CreatePin は新しいピンを作成する
func (s *PinServiceImpl) CreatePin(ctx context.Context, req models.PinCreate, userID string) (*models.Pin, error) {
	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, req.FloorID)
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

	if mapObj.UserID != userID && !mapObj.IsPubliclyEditable {
		return nil, errors.New("このマップにピンを追加する権限がありません")
	}

	// 新しいピンを作成
	pin := &models.Pin{
		FloorID:        req.FloorID,
		Title:          req.Title,
		Description:    req.Description,
		XPosition:      req.XPosition,
		YPosition:      req.YPosition,
		ImageURL:       req.ImageURL,
		EditorID:       req.EditorID,
		EditorNickname: req.EditorNickname,
	}

	// 編集者情報が指定されていない場合、ユーザー情報を使用
	if pin.EditorID == "" {
		pin.EditorID = userID
	}

	if err := s.pinRepo.Create(ctx, pin); err != nil {
		return nil, err
	}

	return pin, nil
}

// GetPinsByFloorID はフロアに属するすべてのピンを取得する
func (s *PinServiceImpl) GetPinsByFloorID(ctx context.Context, floorID string) ([]*models.Pin, error) {
	// フロアIDを使用してピンを取得
	return s.pinRepo.GetByFloorID(ctx, floorID)
}

// GetPinByID はIDによりピンを取得する
func (s *PinServiceImpl) GetPinByID(ctx context.Context, id string) (*models.Pin, error) {
	return s.pinRepo.GetByID(ctx, id)
}

// CreatePin は新しいピンを作成する
func (s *PinServiceImpl) CreatePin(ctx context.Context, req models.PinCreate, userID string) (*models.Pin, error) {
	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, req.FloorID)
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

	if mapObj.UserID != userID && !mapObj.IsPubliclyEditable {
		return nil, errors.New("このマップにピンを追加する権限がありません")
	}

	// 新しいピンを作成
	pin := &models.Pin{
		FloorID:        req.FloorID,
		Title:          req.Title,
		Description:    req.Description,
		XPosition:      req.XPosition,
		YPosition:      req.YPosition,
		ImageURL:       req.ImageURL,
		EditorID:       req.EditorID,
		EditorNickname: req.EditorNickname,
	}

	// 編集者情報が指定されていない場合、ユーザー情報を使用
	if pin.EditorID == "" {
		pin.EditorID = userID
	}

	if err := s.pinRepo.Create(ctx, pin); err != nil {
		return nil, err
	}

	return pin, nil
}

// GetPinsByFloorID はフロアに属するすべてのピンを取得する
func (s *PinServiceImpl) GetPinsByFloorID(ctx context.Context, floorID string) ([]*models.Pin, error) {
	// フロアIDを使用してピンを取得
	return s.pinRepo.GetByFloorID(ctx, floorID)
}

// GetPinByID はIDによりピンを取得する
func (s *PinServiceImpl) GetPinByID(ctx context.Context, id string) (*models.Pin, error) {
	return s.pinRepo.GetByID(ctx, id)
}

// UpdatePin はピン情報を更新する
func (s *PinServiceImpl) UpdatePin(ctx context.Context, id string, req models.PinUpdate, userID string) (*models.Pin, error) {
	// ピンを取得
	pin, err := s.pinRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if pin == nil {
		return nil, errors.New("ピンが見つかりません")
	}

	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, pin.FloorID)
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

	// 権限チェック
	isOwner := mapObj.UserID == userID
	isEditor := pin.EditorID == userID
	isPubliclyEditable := mapObj.IsPubliclyEditable

	if !isOwner && !isEditor && !isPubliclyEditable {
		return nil, errors.New("このピンを編集する権限がありません")
	}

	// ピン情報を更新
	if req.Title != "" {
		pin.Title = req.Title
	}

	if req.Description != "" {
		pin.Description = req.Description
	}

	if req.ImageURL != "" {
		pin.ImageURL = req.ImageURL
	}

	if err := s.pinRepo.Update(ctx, pin); err != nil {
		return nil, err
	}

	return pin, nil
}

// DeletePin はピンを削除する
func (s *PinServiceImpl) DeletePin(ctx context.Context, id string, userID string) error {
	// ピンを取得
	pin, err := s.pinRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if pin == nil {
		return errors.New("ピンが見つかりません")
	}

	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, pin.FloorID)
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

	// 権限チェック
	isOwner := mapObj.UserID == userID
	isEditor := pin.EditorID == userID
	
	if !isOwner && !isEditor {
		return errors.New("このピンを削除する権限がありません")
	}

	// ピンを削除
	return s.pinRepo.Delete(ctx, id)
}

// UploadPinImage はピンの画像をアップロードする
func (s *PinServiceImpl) UploadPinImage(ctx context.Context, pinID string, userID string, file io.Reader) (string, error) {
	// ピンを取得
	pin, err := s.pinRepo.GetByID(ctx, pinID)
	if err != nil {
		return "", err
	}

	if pin == nil {
		return "", errors.New("ピンが見つかりません")
	}

	// フロアを取得
	floor, err := s.floorRepo.GetByID(ctx, pin.FloorID)
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

	// 権限チェック
	isOwner := mapObj.UserID == userID
	isEditor := pin.EditorID == userID
	isPubliclyEditable := mapObj.IsPubliclyEditable

	if !isOwner && !isEditor && !isPubliclyEditable {
		return "", errors.New("このピンを編集する権限がありません")
	}

	// Cloudinaryにアップロード
	uploadResult, err := s.cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "pin_images",
		Format:         "webp",
		Transformation: "q_auto",
		ResourceType:   "image",
		PublicID:       fmt.Sprintf("pin_%s", pinID),
	})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}