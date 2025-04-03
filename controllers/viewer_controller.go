// backend/controllers/viewer_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/services"
)

// ViewerController はビューワー関連の操作を提供するコントローラー
type ViewerController struct {
	mapService   services.MapService
	floorService services.FloorService
	pinService   services.PinService
}

// NewViewerController は新しいViewerControllerを作成する
func NewViewerController(
	mapService services.MapService,
	floorService services.FloorService,
	pinService services.PinService,
) *ViewerController {
	return &ViewerController{
		mapService:   mapService,
		floorService: floorService,
		pinService:   pinService,
	}
}

// GetMapData はマップの全データを取得する
func (c *ViewerController) GetMapData(ctx *gin.Context) {
	mapID := ctx.Param("mapId")
	if mapID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "マップIDが必要です"})
		return
	}

	// マップデータを取得
	mapData, err := c.mapService.GetMapByID(ctx, mapID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの取得に失敗しました"})
		return
	}
	if mapData == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "マップが見つかりません"})
		return
	}

	// フロアデータを取得
	floors, err := c.floorService.GetFloorsByMapID(ctx, mapID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "フロアの取得に失敗しました"})
		return
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
		pins, err = c.pinService.GetByFloorIDs(ctx, floorIDs)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ピンの取得に失敗しました"})
			return
		}
	}

	// レスポンスデータを構築
	responseData := gin.H{
		"map":    mapData,
		"floors": floors,
		"pins":   pins,
	}

	ctx.JSON(http.StatusOK, responseData)
}
