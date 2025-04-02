// backend/controllers/floor_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/mapapp/models"
	"github.com/yourname/mapapp/services"
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

	ctx.JSON(http.StatusOK, floors)
}

// UpdateFloorImage はフロアの画像を更新する
func (c *FloorController) UpdateFloorImage(ctx *gin.Context) {
	floorID := ctx.Param("floorId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// マルチパートフォームからファイルを取得
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルが必要です"})
		return
	}

	// ファイルを開く
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの読み込みに失敗しました"})
		return
	}
	defer src.Close()

	// Cloudinaryにアップロード
	imageURL, err := c.floorService.UploadFloorImage(ctx, floorID, userID.(string), src)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// フロア情報を更新
	floor, err := c.floorService.UpdateFloor(ctx, floorID, models.FloorUpdate{
		ImageURL: imageURL,
	}, userID.(string))
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