// backend/controllers/floor_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/services"
)

// FloorController はフロア関連のAPIエンドポイントを管理する
type FloorController struct {
	floorService services.FloorService
}

// NewFloorController は新しいFloorControllerを作成する
func NewFloorController(floorService services.FloorService) *FloorController {
	return &FloorController{
		floorService: floorService,
	}
}

// CreateFloor は新しいフロアを作成する
func (c *FloorController) CreateFloor(ctx *gin.Context) {
	mapID := ctx.Param("mapId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req models.FloorCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	req.MapID = mapID

	floor, err := c.floorService.CreateFloor(ctx, req, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, floor)
}

// GetFloors はマップに属するすべてのフロアを取得する
func (c *FloorController) GetFloors(ctx *gin.Context) {
	mapID := ctx.Param("mapId")

	floors, err := c.floorService.GetFloorsByMapID(ctx, mapID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if floors == nil {
		floors = []*models.Floor{}
	}

	ctx.JSON(http.StatusOK, floors)
}

// GetFloorByID はIDによりフロアを取得する
func (c *FloorController) GetFloorByID(ctx *gin.Context) {
	floorID := ctx.Param("floorId")

	floor, err := c.floorService.GetFloorByID(ctx, floorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if floor == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "フロアが見つかりません"})
		return
	}

	ctx.JSON(http.StatusOK, floor)
}

// UpdateFloorImage はフロアの画像を更新する
func (c *FloorController) UpdateFloorImage(ctx *gin.Context) {
	floorID := ctx.Param("floorId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// Cloudinaryにアップロード済みの画像URLを取得
	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像URLが必要です"})
		return
	}

	// フロア情報を更新
	update := models.FloorUpdate{
		ImageURL: req.ImageURL,
	}

	floor, err := c.floorService.UpdateFloor(ctx, floorID, update, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, floor)
}

// UpdateFloor はフロア情報を更新する
func (c *FloorController) UpdateFloor(ctx *gin.Context) {
	floorID := ctx.Param("floorId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req models.FloorUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	floor, err := c.floorService.UpdateFloor(ctx, floorID, req, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, floor)
}

// DeleteFloor はフロアを削除する
func (c *FloorController) DeleteFloor(ctx *gin.Context) {
	floorID := ctx.Param("floorId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	err := c.floorService.DeleteFloor(ctx, floorID, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "フロアが正常に削除されました", "id": floorID})
}
