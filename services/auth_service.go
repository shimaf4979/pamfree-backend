// services/auth_service.go
package services

import (
	"context"

	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/repositories"
)

// AuthService 認証サービス
type AuthService struct {
	userRepo repositories.UserRepository
}

// NewAuthService 新しい認証サービスを作成
func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// CreateUser 新しいユーザーを作成
func (s *AuthService) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}

// GetUserByID IDからユーザーを取得
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail メールアドレスからユーザーを取得
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// EmailExists メールアドレスが既に存在するか確認
func (s *AuthService) EmailExists(ctx context.Context, email string) (bool, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// UpdateUser ユーザー情報を更新
func (s *AuthService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser ユーザーを削除
func (s *AuthService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// GetAllUsers すべてのユーザーを取得
func (s *AuthService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.GetAll(ctx)
}
