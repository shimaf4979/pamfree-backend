// backend/controllers/pin_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/services"
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

	pin, err := c.pinService.Create(ctx, userID.(string), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, pin)
}

// GetPinsByFloorID はフロアに属するすべてのピンを取得する
func (c *PinController) GetPinsByFloorID(ctx *gin.Context) {
	floorID := ctx.Param("floorId")

	pins, err := c.pinService.GetByFloorID(ctx, floorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pins == nil {
		pins = []*models.Pin{}
	}

	ctx.JSON(http.StatusOK, pins)
}

// GetPinByID はIDによりピンを取得する
func (c *PinController) GetPinByID(ctx *gin.Context) {
	pinID := ctx.Param("pinId")

	pin, err := c.pinService.GetByID(ctx, pinID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pin == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ピンが見つかりません"})
		return
	}

	ctx.JSON(http.StatusOK, pin)
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

	pin, err := c.pinService.Update(ctx, userID.(string), pinID, &req)
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

	err := c.pinService.Delete(ctx, userID.(string), pinID)
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

	// マルチパートフォームから画像URLを取得
	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像URLが必要です"})
		return
	}

	// Cloudinaryから取得した画像URLでピン情報を更新
	update := &models.PinUpdate{
		ImageURL: req.ImageURL,
	}

	pin, err := c.pinService.Update(ctx, userID.(string), pinID, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pin)
}

// CreatePublicPin は公開編集で新しいピンを作成する
func (c *PinController) CreatePublicPin(ctx *gin.Context) {
	var req models.PinCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	// 必須フィールドの確認
	if req.FloorID == "" || req.Title == "" || req.EditorID == "" || req.EditorNickname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "フロアID、タイトル、編集者情報は必須です"})
		return
	}

	pin, err := c.pinService.CreatePublic(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, pin)
}

// UpdatePublicPin は公開編集でピン情報を更新する
func (c *PinController) UpdatePublicPin(ctx *gin.Context) {
	pinID := ctx.Param("pinId")

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		EditorID    string `json:"editorId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストです"})
		return
	}

	update := &models.PinUpdate{
		Title:       req.Title,
		Description: req.Description,
	}

	pin, err := c.pinService.UpdatePublic(ctx, req.EditorID, pinID, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pin)
}

// DeletePublicPin は公開編集でピンを削除する
func (c *PinController) DeletePublicPin(ctx *gin.Context) {
	pinID := ctx.Param("pinId")
	editorID := ctx.Query("editorId")

	if editorID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "編集者IDは必須です"})
		return
	}

	err := c.pinService.DeletePublic(ctx, editorID, pinID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ピンが正常に削除されました", "id": pinID})
}
