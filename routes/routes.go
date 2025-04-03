// backend/routes/routes.go
package routes

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shimaf4979/pamfree-backend/config"
	"github.com/shimaf4979/pamfree-backend/controllers"
	"github.com/shimaf4979/pamfree-backend/middlewares"
	"github.com/shimaf4979/pamfree-backend/repositories"
	"github.com/shimaf4979/pamfree-backend/services"
)

// SetupRoutes はアプリケーションのルートを設定する
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// データベース接続
	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}

	// リポジトリの初期化
	userRepo := repositories.NewMySQLUserRepository(db)
	mapRepo := repositories.NewMySQLMapRepository(db)
	floorRepo := repositories.NewMySQLFloorRepository(db)
	pinRepo := repositories.NewMySQLPinRepository(db)
	publicEditorRepo := repositories.NewMySQLPublicEditorRepository(db)

	// サービスの初期化
	authService := services.NewAuthService(userRepo)
	mapService := services.NewMapService(mapRepo)
	floorService := services.NewFloorService(floorRepo, mapRepo)
	pinService := services.NewPinService(pinRepo, floorRepo, mapRepo)
	publicEditorService := services.NewPublicEditorService(publicEditorRepo, mapRepo)

	// コントローラーの初期化
	authController := controllers.NewAuthController(authService, cfg.JWTSecret)
	mapController := controllers.NewMapController(mapService)
	floorController := controllers.NewFloorController(floorService)
	pinController := controllers.NewPinController(pinService)
	publicEditorController := controllers.NewPublicEditorController(publicEditorService, mapService)
	viewerController := controllers.NewViewerController(mapService, floorService, pinService)

	// Cloudinaryコントローラー
	cloudinaryController, err := controllers.NewCloudinaryController(cfg)
	if err != nil {
		log.Fatalf("Cloudinaryコントローラーの初期化に失敗しました: %v", err)
	}

	// CORSミドルウェアを設定
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigins},
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		ExposeHeaders:    cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	}
	router.Use(cors.New(corsConfig))

	// 認証ミドルウェア
	authMiddleware := middlewares.AuthMiddleware(cfg.JWTSecret)
	adminMiddleware := middlewares.AdminMiddleware()

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 認証ルート
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/me", authMiddleware, authController.GetMe)
	}

	// マップルート
	maps := router.Group("/api/maps")
	{
		maps.GET("", authMiddleware, mapController.GetMaps)
		maps.POST("", authMiddleware, mapController.CreateMap)
		maps.GET("/:mapId", authMiddleware, mapController.GetMapByID)
		maps.PATCH("/:mapId", authMiddleware, mapController.UpdateMap)
		maps.DELETE("/:mapId", authMiddleware, mapController.DeleteMap)

		// フロアルート (マップIDによる)
		maps.GET("/:mapId/floors", floorController.GetFloors)
		maps.POST("/:mapId/floors", authMiddleware, floorController.CreateFloor)
	}

	// フロアルート
	floors := router.Group("/api/floors")
	{
		floors.GET("/:floorId", floorController.GetFloorByID)
		floors.PATCH("/:floorId", authMiddleware, floorController.UpdateFloor)
		floors.DELETE("/:floorId", authMiddleware, floorController.DeleteFloor)

		// ピンルート (フロアIDによる)
		floors.GET("/:floorId/pins", pinController.GetPinsByFloorID)
		floors.POST("/:floorId/pins", authMiddleware, pinController.CreatePin)

		// フロア画像アップロード
		floors.POST("/:floorId/image", authMiddleware, floorController.UpdateFloorImage)
	}

	// ピンルート
	pins := router.Group("/api/pins")
	{
		pins.GET("/:pinId", pinController.GetPinByID)
		pins.PATCH("/:pinId", authMiddleware, pinController.UpdatePin)
		pins.DELETE("/:pinId", authMiddleware, pinController.DeletePin)

		// ピン画像アップロード
		pins.POST("/:pinId/image", authMiddleware, pinController.UpdatePinImage)
	}

	// 公開編集ルート
	publicEdit := router.Group("/api/public-edit")
	{
		publicEdit.POST("/register", publicEditorController.Register)
		publicEdit.POST("/verify", publicEditorController.Verify)

		// 公開編集用のピン操作
		publicEdit.POST("/pins", pinController.CreatePublicPin)
		publicEdit.PATCH("/pins/:pinId", pinController.UpdatePublicPin)
		publicEdit.DELETE("/pins/:pinId", pinController.DeletePublicPin)
	}

	// ビューワールート
	viewer := router.Group("/api/viewer")
	{
		viewer.GET("/:mapId", viewerController.GetMapData)
	}

	// Cloudinaryルート
	cloudinary := router.Group("/api/cloudinary", authMiddleware)
	{
		cloudinary.POST("/upload", cloudinaryController.UploadImage)
		cloudinary.POST("/delete", cloudinaryController.DeleteImage)
	}

	// アカウント管理ルート
	account := router.Group("/api/account", authMiddleware)
	{
		account.PATCH("/update-profile", authController.UpdateProfile)
		account.POST("/change-password", authController.ChangePassword)
	}

	// 管理者ルート
	admin := router.Group("/api/admin", authMiddleware, adminMiddleware)
	{
		admin.GET("/users", authController.GetAllUsers)
		admin.PATCH("/users/:userId", authController.UpdateUser)
		admin.DELETE("/users/:userId", authController.DeleteUser)
	}
}

// setupDatabase はデータベース接続を設定する
func setupDatabase(cfg *config.Config) (*sql.DB, error) {
	// データソース名を構築
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// データベース接続を開く
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 接続をテスト
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 接続プールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
