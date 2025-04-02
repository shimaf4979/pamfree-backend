// backend/middlewares/auth_middleware.go
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/utils"
)

// AuthMiddleware は認証を検証するミドルウェア
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			c.Abort()
			return
		}

		// Bearerトークンを取得
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証形式が不正です"})
			c.Abort()
			return
		}

		// トークンを検証
		claims, err := utils.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		// ユーザー情報をコンテキストに保存
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// AdminMiddleware は管理者権限を検証するミドルウェア
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "管理者権限が必要です"})
			c.Abort()
			return
		}
		c.Next()
	}
}
