package router

import (
	"fmt"
	"os"
	"time"

	"w3st/interfaces/middlewares"

	"w3st/factory"
	"w3st/infra"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {
	// 環境変数 PORT を取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // デフォルトポート
	}

	db := infra.SetupDB()
	f := factory.NewFactory(db)

	// ルーターの初期化
	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},                                                                                                           // 許可するオリジン
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                                                                     // 許可するHTTPメソッド
		AllowHeaders: []string{"Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Origin", "Content-Type", "Authorization"}, // 許可するヘッダー
		MaxAge:       12 * time.Hour,                                                                                                          // キャッシュの最大時間
	}))

	// connectionTest
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "connection success!!!!",
		})
	})

	// ユーザー
	users := r.Group("/users")
	userController := f.InitUserController()

	// Auth
	authUsecase := f.InitAuthUsecase()

	// API Key Auth
	apiKeyUsecase := f.InitApiKeyUsecase()
	apiKeyController := f.InitApiKeyController()

	// Collections - SDK専用 (APIキー認証 + プロジェクトレート制限)
	sdkCollections := r.Group("/collections")
	sdkCollections.Use(middlewares.ApiKeyAuthMiddleware(apiKeyUsecase))
	projectUsecase := f.InitProjectUsecase()
	systemAlertUsecase := f.InitSystemAlertUsecase()
	sdkCollections.Use(middlewares.ProjectRateLimitMiddleware(projectUsecase, systemAlertUsecase))
	sdkCollectionController := f.InitSDKCollectionsController()
	sdkEntriesController := f.InitSDKEntriesController()

	// API - GUI専用 (Auth0認証)
	api := r.Group("/api")
	api.Use(middlewares.Auth0AuthMiddleware())
	guiCollectionController := f.InitGUICollectionsController()
	guiEntriesController := f.InitGUIEntriesController()

	// Media
	mediaController := f.InitMediaController()

	// Versions
	versionController := f.InitVersionController()

	// Permissions
	permissionController := f.InitPermissionController()

	// Audit
	auditController := f.InitAuditController()

	// System Alerts
	systemAlertController := f.InitSystemAlertController()

	// Projects
	projectController := f.InitProjectController()

	// ユーザー登録
	users.POST("/signup", userController.Signup)

	// ログイン
	users.POST("/login", userController.Login)

	// ユーザー情報取得
	users.GET("/me", middlewares.JwtAuthMiddleware(authUsecase), userController.GetUserInfo)
	// ユーザー情報更新
	users.PUT("/me", middlewares.JwtAuthMiddleware(authUsecase), userController.UpdateUser)

	// 管理者用ユーザー管理API
	apiUsers := api.Group("/users")
	apiUsers.GET("", userController.GetAllUsers)
	apiUsers.PUT("/:userId", userController.UpdateUserById)
	apiUsers.DELETE("/:userId", userController.DeleteUser)

	// API Keys - GUI専用
	api.POST("/api-keys", apiKeyController.CreateApiKey)

	// SDK専用ルート
	// Collection一覧を取得
	sdkCollections.GET("", sdkCollectionController.GetCollectionByProjectId)
	// Collection詳細取得
	sdkCollections.GET("/:collectionId", sdkCollectionController.GetCollectionsByCollectionId)
	// Entries - SDK専用
	sdkEntries := sdkCollections.Group("/:collectionId/entries")
	sdkEntries.GET("", sdkEntriesController.GetEntries)

	// GUI専用ルート
	// Collection一覧取得
	api.GET("/collections", guiCollectionController.GetCollections)
	// Collection作成
	api.POST("/collections", guiCollectionController.MakeCollection)
	// Collection更新
	api.PUT("/collections/:collectionId", guiCollectionController.UpdateCollection)
	// Collection削除
	api.DELETE("/collections/:collectionId", guiCollectionController.DeleteCollection)
	// Fields
	apiFields := api.Group("/collections/:collectionId/fields")
	apiFields.GET("", guiCollectionController.GetFields)
	apiFields.POST("", guiCollectionController.CreateField)
	apiFields.PUT("/:fieldId", guiCollectionController.UpdateField)
	apiFields.DELETE("/:fieldId", guiCollectionController.DeleteField)

	// Entries - GUI専用
	guiEntries := api.Group("/collections/:collectionId/entries")
	guiEntries.GET("", guiEntriesController.GetEntries)
	guiEntries.POST("", guiEntriesController.CreateEntry)
	guiEntries.PUT("/:entryId", guiEntriesController.UpdateEntry)
	guiEntries.DELETE("/:entryId", guiEntriesController.DeleteEntry)

	// Media routes - GUI専用
	api.POST("/media", mediaController.Upload)
	api.GET("/media", mediaController.GetByUserID)
	api.GET("/media/:id", mediaController.GetByID)
	api.DELETE("/media/:id", mediaController.Delete)

	// Versions routes - GUI専用
	api.POST("/versions", versionController.CreateVersion)
	api.GET("/versions/:contentID", versionController.GetVersionsByContentID)
	api.GET("/versions/:contentID/latest", versionController.GetLatestVersion)
	api.POST("/versions/:contentID/restore/:versionID", versionController.RestoreVersion)

	// Permissions routes - GUI専用
	api.GET("/permissions/check", permissionController.CheckPermission)
	api.POST("/permissions/grant", permissionController.GrantPermission)
	api.POST("/permissions/revoke", permissionController.RevokePermission)
	api.GET("/permissions/user", permissionController.GetPermissionsByUser)

	// Audit routes - GUI専用
	api.POST("/audit", auditController.LogAction)
	api.GET("/audit/user", auditController.GetLogsByUser)
	api.GET("/audit/action/:action", auditController.GetLogsByAction)
	api.GET("/audit/project/:projectId", auditController.GetLogsByProject)
	api.GET("/audit/all", auditController.GetAllLogs)

	// System Alert routes - GUI専用
	api.GET("/system-alerts", systemAlertController.GetAlerts)
	api.GET("/system-alerts/active", systemAlertController.GetActiveAlerts)
	api.POST("/system-alerts", systemAlertController.CreateAlert)
	api.PUT("/system-alerts/:id/read", systemAlertController.MarkAsRead)
	api.DELETE("/system-alerts/:id", systemAlertController.DeleteAlert)
	api.GET("/system-alerts/count", systemAlertController.GetAlertCount)

	// Project routes - GUI専用
	api.POST("/projects", projectController.CreateProject)
	api.GET("/projects", projectController.GetAllProjects)
	api.GET("/projects/:id", projectController.GetProjectByID)

	// 指定されたポートでサーバーを開始
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
