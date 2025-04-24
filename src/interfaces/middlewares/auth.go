package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"w3st/usecase"

	myerrors "w3st/errors"
	"w3st/interfaces/controllers"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(authUsecase usecase.JwtUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// tokenをヘッダーから取得
		authHeader := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// tokenの存在を確認
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		//　tokenの検証
		userID, err := authUsecase.ValidateToken(token)
		if err != nil {
			domainErr := &myerrors.DomainError{}
			if errors.As(err, &domainErr) {
				err := controllers.ErrorHandle(domainErr)
				c.JSON(controllers.HttpStatusCodeFromConnectCode(err.Code()), gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}
		// tokenの検証に成功した場合、userIDをコンテキストに保存
		c.Set("userID", userID)

		// tokenが有効な場合、次のハンドラーに進む
		c.Next()
	}
}
