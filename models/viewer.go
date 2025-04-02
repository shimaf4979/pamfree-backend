// backend/models/viewer.go
package models

// ViewerData はビューワー向けのデータを表す構造体
type ViewerData struct {
	Map    *Map     `json:"map"`
	Floors []*Floor `json:"floors"`
	Pins   []*Pin   `json:"pins"`
}
