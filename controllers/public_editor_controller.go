// backend/controllers/public_editor_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/models"
	"github.com/shimaf4979/pamfree-backend/services"
)

// PublicEditorController は公開編集者関連の操作を提供するコントローラー
type PublicEditorController struct {
	publicEditorService services.PublicEditorService
	mapService          services.MapService
}

// NewPublicEditorController は新しいPublicEditorControllerを作成する
func NewPublicEditorController(
	publicEditorService services.PublicEditorService,
	mapService services.MapService,
) *PublicEditorController {
	return &PublicEditorController{
		publicEditorService: publicEditorService,
		mapService:          mapService,
	}
}

// Register は新しい公開編集者を登録する
func (c *PublicEditorController) Register(ctx *gin.Context) {
	var req models.PublicEditorRegister
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	// マップが存在するか確認
	mapData, err := c.mapService.GetMapByMapID(ctx, req.MapID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "マップの取得に失敗しました"})
		return
	}
	if mapData == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "マップが見つかりません"})
		return
	}

	// マップが公開編集可能か確認
	if !mapData.IsPubliclyEditable {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "このマップは公開編集が許可されていません"})
		return
	}

	// 編集者を登録
	editor, err := c.publicEditorService.Register(ctx, req.MapID, req.Nickname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "編集者の登録に失敗しました"})
		return
	}

	// レスポンスを構築
	response := models.PublicEditorResponse{
		EditorID: editor.ID,
		Nickname: editor.Nickname,
		Token:    editor.EditorToken,
		MapID:    editor.MapID,
		Verified: true,
	}

	ctx.JSON(http.StatusCreated, response)
}

// Verify は公開編集者トークンを検証する
func (c *PublicEditorController) Verify(ctx *gin.Context) {
	var req models.PublicEditorVerify
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	// トークンを検証
	editor, err := c.publicEditorService.Verify(ctx, req.EditorID, req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです", "verified": false})
		return
	}

	// 最終アクティブ時間を更新
	if err := c.publicEditorService.UpdateLastActive(ctx, req.EditorID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新に失敗しました", "verified": false})
		return
	}

	// レスポンスを構築
	response := models.PublicEditorResponse{
		EditorID: editor.ID,
		Nickname: editor.Nickname,
		MapID:    editor.MapID,
		Verified: true,
	}

	ctx.JSON(http.StatusOK, response)
}
