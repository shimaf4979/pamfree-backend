// backend/routes/routes.go
package routes

import (
	"database/sql"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yourname/mapapp/config"
	"github.com/yourname/mapapp/controllers"
	"github.com/yourname/mapapp/middlewares"
	"github.com/yourname/mapapp/repositories"
	"github.com/yourname/mapapp/services"
)

// SetupRoutes はルートの設定を行う
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// データベース接続
	dsn := getDSN(cfg)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("データベース接続に失敗しました: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic("データベース接続の確認に失敗しました: " + err.Error())
	}

	// Cloudinaryの初期化
	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryName,
		cfg.CloudinaryKey,
		cfg.CloudinarySecret,
	)
	if err != nil {
		panic("Cloudinaryの初期化に失敗しました: " + err.Error())
	}

	// リポジトリの初期化
	userRepo := repositories.NewMySQLUserRepository(db)
	mapRepo := repositories.NewMySQLMapRepository(db)
	floorRepo := repositories.NewMySQLFloorRepository(db)
	pinRepo := repositories.NewMySQLPinRepository(db)
	publicEditorRepo := repositories.NewMySQLPublicEditorRepository(db)

	// サービスの初期化
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	mapService := services.NewMapService(mapRepo)
	floorService := services.NewFloorService(floorRepo, mapRepo, cld)
	pinService := services.NewPinService(pinRepo, floorRepo, mapRepo, cld)
	publicEditorService := services.NewPublicEditorService(publicEditorRepo, mapRepo)
	viewerService := services.NewViewerService(mapRepo, floorRepo, pinRepo)

	// コントローラーの初期化
	authController := controllers.NewAuthController(authService)
	mapController := controllers.NewMapController(mapService)
	floorController := controllers.NewFloorController(floorService)
	pinController := controllers.NewPinController(pinService)
	publicEditorController := controllers.NewPublicEditorController(publicEditorService)
	viewerController := controllers.NewViewerController(viewerService)

	// CORSミドルウェア
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 認証ミドルウェア
	authMiddleware := middlewares.AuthMiddleware(cfg.JWTSecret)
	adminMiddleware := middlewares.AdminMiddleware()

	// 認証関連のルート
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// マップ関連のルート
	maps := router.Group("/maps")
	maps.Use(authMiddleware)
	{
		maps.GET("", mapController.GetMaps)
		maps.POST("", mapController.CreateMap)
		maps.GET("/:mapId", mapController.GetMap)
		maps.PATCH("/:mapId", mapController.UpdateMap)
		maps.DELETE("/:mapId", mapController.DeleteMap)

		// フロア関連のルート
		maps.GET("/:mapId/floors", floorController.GetFloors)
		maps.POST("/:mapId/floors", floorController.CreateFloor)
	}

	// フロア関連のルート
	floors := router.Group("/floors")
	floors.Use(authMiddleware)
	{
		floors.PATCH("/:floorId", floorController.UpdateFloor)
		floors.DELETE("/:floorId", floorController.DeleteFloor)
		floors.POST("/:floorId/image", floorController.UpdateFloorImage)

		// ピン関連のルート
		floors.GET("/:floorId/pins", pinController.GetPins)
		floors.POST("/:floorId/pins", pinController.CreatePin)
	}

	// ピン関連のルート
	pins := router.Group("/pins")
	pins.Use(authMiddleware)
	{
		pins.PATCH("/:pinId", pinController.UpdatePin)
		pins.DELETE("/:pinId", pinController.DeletePin)
		pins.POST("/:pinId/image", pinController.UpdatePinImage)
	}

	// 管理者用ルート
	admin := router.Group("/admin")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.GET("/users", authController.GetUsers)
		admin.PATCH("/users/:userId", authController.UpdateUser)
		admin.DELETE("/users/:userId", authController.DeleteUser)
	}

	// 公開閲覧用ルート
	viewer := router.Group("/viewer")
	{
		viewer.GET("/:mapId", viewerController.GetMapData)
	}

	// 公開編集用ルート
	publicEdit := router.Group("/public-edit")
	{
		publicEdit.POST("/register", publicEditorController.Register)
		publicEdit.POST("/verify", publicEditorController.VerifyToken)
		publicEdit.POST("/pins", publicEditorController.CreatePin)
		publicEdit.PATCH("/pins/:pinId", publicEditorController.UpdatePin)
		publicEdit.DELETE("/pins/:pinId", publicEditorController.DeletePin)
	}
}

// getDSN はデータベース接続文字列を生成する
func getDSN(cfg *config.Config) string {
	return cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?parseTime=true"
}