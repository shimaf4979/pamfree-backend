// backend/controllers/cloudinary_controller.go
package controllers

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/shimaf4979/pamfree-backend/config"
)

// CloudinaryController はCloudinaryとの画像連携を行うコントローラー
type CloudinaryController struct {
	cloudinary *cloudinary.Cloudinary
	config     *config.Config
}

// NewCloudinaryController は新しいCloudinaryControllerを作成する
func NewCloudinaryController(cfg *config.Config) (*CloudinaryController, error) {
	// Cloudinaryクライアントの初期化
	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryName,
		cfg.CloudinaryKey,
		cfg.CloudinarySecret,
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryController{
		cloudinary: cld,
		config:     cfg,
	}, nil
}

// UploadImage は画像をCloudinaryにアップロードする
func (c *CloudinaryController) UploadImage(ctx *gin.Context) {
	// マルチパートフォームファイルを取得
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルが必要です"})
		return
	}
	defer file.Close()

	// ファイルタイプの検証
	if !isValidImageType(header.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効な画像形式です"})
		return
	}

	// アップロードフォルダ（任意）
	folder := ctx.DefaultQuery("folder", "images")

	// アップロードオプション
	uploadParams := uploader.UploadParams{
		Folder:         folder,
		Transformation: "f_auto,q_auto", // 自動フォーマットと品質最適化
	}

	// Cloudinaryにアップロード
	uploadResult, err := c.cloudinary.Upload.Upload(
		context.Background(),
		file,
		uploadParams,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "画像のアップロードに失敗しました"})
		return
	}

	// レスポンスを返す
	ctx.JSON(http.StatusOK, gin.H{
		"url":        uploadResult.SecureURL,
		"public_id":  uploadResult.PublicID,
		"format":     uploadResult.Format,
		"width":      uploadResult.Width,
		"height":     uploadResult.Height,
		"bytes":      uploadResult.Bytes,
		"created_at": uploadResult.CreatedAt,
	})
}

// DeleteImage はCloudinaryから画像を削除する
func (c *CloudinaryController) DeleteImage(ctx *gin.Context) {
	var req struct {
		PublicID string `json:"publicId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "公開IDが必要です"})
		return
	}

	// Cloudinaryから画像を削除
	result, err := c.cloudinary.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{PublicID: req.PublicID},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "画像の削除に失敗しました"})
		return
	}

	if result.Result != "ok" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "画像の削除に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "画像が正常に削除されました"})
}

// isValidImageType は有効な画像形式かどうかをチェックする
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
		".tiff": true,
		".svg":  true,
	}
	return validExts[ext]
}
