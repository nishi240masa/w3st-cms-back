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

	// Collections
	collections := r.Group("/collections")
	collectionController := f.InitCollectionController()
	// Fields
	fieldController := f.InitFieldController()

	// ユーザー登録
	users.POST("/signup", userController.Signup)

	// ログイン
	users.POST("/login", userController.Login)

	// ユーザー情報取得
	users.GET("/me", middlewares.JwtAuthMiddleware(authUsecase), userController.GetUserInfo)

	// Collectionを追加
	collections.POST("", middlewares.JwtAuthMiddleware(authUsecase), collectionController.MakeCollection)

	// Collection一覧を取得
	collections.GET("", middlewares.JwtAuthMiddleware(authUsecase), collectionController.GetCollectionByUserId)

	// Collectionを取得
	// Fields
	fields := collections.Group("/:collectionId/fields")
	fields.POST("", middlewares.JwtAuthMiddleware(authUsecase), fieldController.Create)
	fields.PUT("/:fieldId", middlewares.JwtAuthMiddleware(authUsecase), fieldController.Update)
	fields.DELETE("/:fieldId", middlewares.JwtAuthMiddleware(authUsecase), fieldController.Delete)
	collections.GET("/:collectionId", middlewares.JwtAuthMiddleware(authUsecase), collectionController.GetCollectionsByCollectionId)

	// 指定されたポートでサーバーを開始
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
