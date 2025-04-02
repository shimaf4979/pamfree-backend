// backend/models/pin.go
package models

import (
	"time"
)

// Pin はピン情報を表す構造体
type Pin struct {
	ID             string    `json:"id" db:"id"`
	FloorID        string    `json:"floor_id" db:"floor_id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	XPosition      float64   `json:"x_position" db:"x_position"`
	YPosition      float64   `json:"y_position" db:"y_position"`
	ImageURL       string    `json:"image_url" db:"image_url"`
	EditorID       string    `json:"editor_id" db:"editor_id"`
	EditorNickname string    `json:"editor_nickname" db:"editor_nickname"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// PinCreate はピン作成リクエストを表す構造体
type PinCreate struct {
	FloorID        string  `json:"floor_id" binding:"required"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	XPosition      float64 `json:"x_position" binding:"required"`
	YPosition      float64 `json:"y_position" binding:"required"`
	ImageURL       string  `json:"image_url"`
	EditorID       string  `json:"editor_id"`
	EditorNickname string  `json:"editor_nickname"`
}

// PinUpdate はピン更新リクエストを表す構造体
type PinUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
