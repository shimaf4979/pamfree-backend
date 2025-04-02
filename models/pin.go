// backend/models/floor.go
package models

import (
	"time"
)

// Floor はフロア（エリア）情報を表す構造体
type Floor struct {
	ID          string    `json:"id" db:"id"`
	MapID       string    `json:"map_id" db:"map_id"`
	FloorNumber int       `json:"floor_number" db:"floor_number"`
	Name        string    `json:"name" db:"name"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// FloorCreate はフロア作成リクエストを表す構造体
type FloorCreate struct {
	MapID       string `json:"map_id" binding:"required"`
	FloorNumber int    `json:"floor_number" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

// FloorUpdate はフロア更新リクエストを表す構造体
type FloorUpdate struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}