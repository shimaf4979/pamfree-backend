// backend/models/map.go
package models

import (
	"time"
)

// Map はマップ情報を表す構造体
type Map struct {
	ID                 string    `json:"id" db:"id"`
	Title              string    `json:"title" db:"title"`
	Description        string    `json:"description" db:"description"`
	UserID             string    `json:"user_id" db:"user_id"`
	IsPubliclyEditable bool      `json:"is_publicly_editable" db:"is_publicly_editable"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// MapCreate はマップ作成リクエストを表す構造体
type MapCreate struct {
	ID                 string `json:"id" binding:"required"`
	Title              string `json:"title" binding:"required"`
	Description        string `json:"description"`
	IsPubliclyEditable bool   `json:"is_publicly_editable"`
}

// MapUpdate はマップ更新リクエストを表す構造体
type MapUpdate struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	IsPubliclyEditable bool   `json:"is_publicly_editable"`
}
