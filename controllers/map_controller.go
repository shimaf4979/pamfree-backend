// controllers/map_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/services"
)

// MapController マップコントローラー
type MapController struct {
	mapService services.MapService
}

// NewMapController 新しいマップコントローラーを作成
func NewMapController(mapService services.MapService) *MapController {
	return &MapController{
		mapService: mapService,
	}
}

// GetMaps ユーザーのマップ一覧取得ハンドラー
func (c *MapController) GetMaps(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	maps, err := c.mapService.GetMapsByUserID(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの取得に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, maps)
}

// GetMapByID マップ取得ハンドラー
func (c *MapController) GetMapByID(ctx *gin.Context) {
	mapID := ctx.Param("mapId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// マップを取得
	m, err := c.mapService.GetMapByID(ctx, mapID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの取得に失敗しました"})
		return
	}
	if m == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "マップが見つかりません"})
		return
	}

	// マップの所有者を確認
	if m.UserID != userID.(string) {
		userRole, exists := ctx.Get("userRole")
		if !exists || userRole.(string) != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "このマップにアクセスする権限がありません"})
			return
		}
	}

	ctx.JSON(http.StatusOK, m)
}

// CreateMap マップ作成ハンドラー
func (c *MapController) CreateMap(ctx *gin.Context) {
	var req models.MapCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// マップの作成
	m := &models.Map{
		ID:                 req.ID, // クライアントからIDを受け取る
		Title:              req.Title,
		Description:        req.Description,
		UserID:             userID.(string),
		IsPubliclyEditable: req.IsPubliclyEditable,
	}

	if err := c.mapService.CreateMap(ctx, m); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの作成に失敗しました"})
		return
	}

	ctx.JSON(http.StatusCreated, m)
}

// UpdateMap マップ更新ハンドラー
func (c *MapController) UpdateMap(ctx *gin.Context) {
	mapID := ctx.Param("mapId")
	var req models.MapUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// マップを取得
	m, err := c.mapService.GetMapByID(ctx, mapID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの取得に失敗しました"})
		return
	}
	if m == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "マップが見つかりません"})
		return
	}

	// マップの所有者を確認
	if m.UserID != userID.(string) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "このマップを編集する権限がありません"})
		return
	}

	// マップを更新
	if req.Title != "" {
		m.Title = req.Title
	}
	m.Description = req.Description
	m.IsPubliclyEditable = req.IsPubliclyEditable

	if err := c.mapService.UpdateMap(ctx, m); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの更新に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, m)
}

// DeleteMap マップ削除ハンドラー
func (c *MapController) DeleteMap(ctx *gin.Context) {
	mapID := ctx.Param("mapId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// マップを取得
	m, err := c.mapService.GetMapByID(ctx, mapID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの取得に失敗しました"})
		return
	}
	if m == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "マップが見つかりません"})
		return
	}

	// マップの所有者または管理者を確認
	if m.UserID != userID.(string) {
		userRole, exists := ctx.Get("userRole")
		if !exists || userRole.(string) != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "このマップを削除する権限がありません"})
			return
		}
	}

	if err := c.mapService.DeleteMap(ctx, mapID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの削除に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "マップが正常に削除されました"})
}
