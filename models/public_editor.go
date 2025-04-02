// backend/models/public_editor.go
package models

import (
	"time"
)

// PublicEditor は公開編集者情報を表す構造体
type PublicEditor struct {
	ID          string    `json:"id" db:"id"`
	MapID       string    `json:"map_id" db:"map_id"`
	Nickname    string    `json:"nickname" db:"nickname"`
	EditorToken string    `json:"-" db:"editor_token"` // トークンはJSONに含めない
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LastActive  time.Time `json:"last_active" db:"last_active"`
}

// PublicEditorRegister は公開編集者登録リクエストを表す構造体
type PublicEditorRegister struct {
	MapID    string `json:"mapId" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// PublicEditorVerify は公開編集者検証リクエストを表す構造体
type PublicEditorVerify struct {
	EditorID string `json:"editorId" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

// PublicEditorResponse は公開編集者登録レスポンスを表す構造体
type PublicEditorResponse struct {
	EditorID string `json:"editorId"`
	Nickname string `json:"nickname"`
	Token    string `json:"token,omitempty"` // 登録時のみトークンを含める
	MapID    string `json:"mapId"`
	Verified bool   `json:"verified"`
}

// ToResponse は公開編集者モデルからレスポンスモデルに変換する
func (e *PublicEditor) ToResponse(includeToken bool) PublicEditorResponse {
	response := PublicEditorResponse{
		EditorID: e.ID,
		Nickname: e.Nickname,
		MapID:    e.MapID,
		Verified: true,
	}

	if includeToken {
		response.Token = e.EditorToken
	}

	return response
}
