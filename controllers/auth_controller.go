// controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/mapapp/models"
	"github.com/yourname/mapapp/services"
	"github.com/yourname/mapapp/utils"
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