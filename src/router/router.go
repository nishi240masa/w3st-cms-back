package router

import (
	"fmt"
	"os"
	"time"
	"w3st/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {
	// 環境変数 PORT を取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // ここで特定のオリジンを許可することもできます（例: []string{"http://localhost:3000"})
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 許可するHTTPメソッド
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 許可するヘッダー
		AllowCredentials: true,                                                // クレデンシャルを許可するかどうか
		MaxAge:           12 * time.Hour,                                      // キャッシュの最大時間
	}))

	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "connection success!!!!",
		})

	})

	// ユーザー
	users := r.Group("/users")
	// 全部取得
	// users.GET("/", controllers.GetUsers)

	// ユーザー登録
	users.POST("/signup", controllers.Signup)

	// ログイン
	users.POST("/login", controllers.Login)


	// サイト

	// 作品
	// productions := r.Group("/productions")
	// 全部取得
	// productions.GET("/", controllers.GetProductions)
	// 特定の作品取得
	// productions.GET("/:id", controllers.GetProduction)



	// 指定されたポートでサーバーを開始
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
