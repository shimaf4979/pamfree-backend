// backend/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourname/mapapp/config"
	"github.com/yourname/mapapp/routes"
)

func main() {
	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// Ginモードの設定
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Ginルーターの初期化
	router := gin.Default()

	// ルートの設定
	routes.SetupRoutes(router, cfg)

	// サーバー起動
	log.Printf("サーバーを起動します: %s\n", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}