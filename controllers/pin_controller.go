// backend/controllers/pin_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/mapapp/models"
	"github.com/yourname/mapapp/services"
)

// PinController はピン関連のAPIエンドポイントを管理する
type PinController struct {
	pinService services.PinService
}

// NewPinController は新しいPinControllerを作成する
func NewPinController(pinService services.PinService) *PinController {
	return &PinController{
		pinService: pinService,
	}
}

// CreatePin は新しいピンを作成する
func (c *PinController) CreatePin(ctx *gin.Context) {
	floorID := ctx.Param("floorId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req models.PinCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	req.FloorID = floorID

	pin, err := c.pinService.CreatePin(ctx, req, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, pin)
}

// GetPins はフロアに属するすべてのピンを取得する
func (c *PinController) GetPins(ctx *gin.Context) {
	floorID := ctx.Param("floorId")

	pins, err := c.pinService.GetPinsByFloorID(ctx, floorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pins)
}

// UpdatePin はピン情報を更新する
func (c *PinController) UpdatePin(ctx *gin.Context) {
	pinID := ctx.Param("pinId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	var req models.PinUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	pin, err := c.pinService.UpdatePin(ctx, pinID, req, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pin)
}

// DeletePin はピンを削除する
func (c *PinController) DeletePin(ctx *gin.Context) {
	pinID := ctx.Param("pinId")
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	err := c.pinService.DeletePin(ctx, pinID, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ピンが正常に削除されました", "id": pinID})
}

// UpdatePinImage はピンの画像を更新する
func (c *PinController) UpdatePinImage(ctx *gin.Context) {
	pinID := ctx.Param("pinId")
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
	imageURL, err := c.pinService.UploadPinImage(ctx, pinID, userID.(string), src)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ピン情報を更新
	pin, err := c.pinService.UpdatePin(ctx, pinID, models.PinUpdate{
		ImageURL: imageURL,
	}, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pin)
}