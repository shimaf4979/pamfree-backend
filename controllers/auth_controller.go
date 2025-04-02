// controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/services"
	"github.com/shimaf4979/pamfree-backend/utils"
)

// AuthController 認証コントローラー
type AuthController struct {
	authService *services.AuthService
	jwtSecret   string
}

// NewAuthController 新しい認証コントローラーを作成
func NewAuthController(authService *services.AuthService, jwtSecret string) *AuthController {
	return &AuthController{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}

// Register ユーザー登録ハンドラー
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.UserRegister
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	// メールアドレスの重複チェック
	exists, err := c.authService.EmailExists(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "このメールアドレスは既に登録されています"})
		return
	}

	// パスワードのハッシュ化
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの処理に失敗しました"})
		return
	}

	// ユーザーの作成
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     "user", // デフォルトは一般ユーザー
	}

	if err := c.authService.CreateUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーの登録に失敗しました"})
		return
	}

	// レスポンスから機密情報を削除
	userResponse := user.ToResponse()

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "ユーザーが正常に登録されました",
		"user":    userResponse,
	})
}

// Login ログインハンドラー
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.UserLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	// メールアドレスからユーザーを検索
	user, err := c.authService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}

	// パスワードの検証
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}

	// JWTトークンの生成
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, c.jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "認証トークンの生成に失敗しました"})
		return
	}

	// レスポンスから機密情報を削除
	userResponse := user.ToResponse()

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  userResponse,
	})
}

// GetMe 現在のユーザー情報取得ハンドラー
func (c *AuthController) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// ユーザーIDからユーザーを検索
	user, err := c.authService.GetUserByID(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	// レスポンスから機密情報を削除
	userResponse := user.ToResponse()

	ctx.JSON(http.StatusOK, gin.H{"user": userResponse})
}

// UpdateProfile ユーザープロフィール更新ハンドラー
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	// ユーザー取得
	user, err := c.authService.GetUserByID(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	// 名前を更新
	user.Name = req.Name

	// ユーザー情報を更新
	if err := c.authService.UpdateUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "プロフィールの更新に失敗しました"})
		return
	}

	// レスポンスから機密情報を削除
	userResponse := user.ToResponse()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "プロフィールを更新しました",
		"user":    userResponse,
	})
}

// ChangePassword パスワード変更ハンドラー
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	// ユーザー取得
	user, err := c.authService.GetUserByID(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	// 現在のパスワードを検証
	if err := utils.CheckPassword(user.Password, req.CurrentPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "現在のパスワードが正しくありません"})
		return
	}

	// 新しいパスワードをハッシュ化
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの処理に失敗しました"})
		return
	}

	// パスワードを更新
	user.Password = hashedPassword

	// ユーザー情報を更新
	if err := c.authService.UpdateUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの更新に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "パスワードを更新しました",
	})
}

// GetAllUsers すべてのユーザーを取得（管理者用）
func (c *AuthController) GetAllUsers(ctx *gin.Context) {
	users, err := c.authService.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー一覧の取得に失敗しました"})
		return
	}

	// レスポンスから機密情報を削除
	var usersResponse []models.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, user.ToResponse())
	}

	ctx.JSON(http.StatusOK, usersResponse)
}

// UpdateUser 特定のユーザーを更新（管理者用）
func (c *AuthController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("userId")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	// ロールのバリデーション
	if req.Role != "admin" && req.Role != "user" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "有効な役割を指定してください"})
		return
	}

	// ユーザー取得
	user, err := c.authService.GetUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	// 自分自身のロールは変更不可
	adminID, exists := ctx.Get("userID")
	if exists && adminID.(string) == userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "自分自身の役割は変更できません"})
		return
	}

	// ロールを更新
	user.Role = req.Role

	// ユーザー情報を更新
	if err := c.authService.UpdateUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーの更新に失敗しました"})
		return
	}

	// レスポンスから機密情報を削除
	userResponse := user.ToResponse()

	ctx.JSON(http.StatusOK, userResponse)
}

// DeleteUser 特定のユーザーを削除（管理者用）
func (c *AuthController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("userId")

	// 自分自身の削除は不可
	adminID, exists := ctx.Get("userID")
	if exists && adminID.(string) == userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "自分自身を削除することはできません"})
		return
	}

	// ユーザー取得
	user, err := c.authService.GetUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	// ユーザーを削除
	if err := c.authService.DeleteUser(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーの削除に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ユーザーが正常に削除されました",
	})
}
