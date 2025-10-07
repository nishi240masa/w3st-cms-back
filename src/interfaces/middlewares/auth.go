package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"

	myerrors "w3st/errors"
	"w3st/interfaces/controllers"
	"w3st/usecase"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApiKeyClaims struct {
	UserID        uuid.UUID `json:"user_id"`
	ProjectID     int       `json:"project_id"`
	CollectionIds []int     `json:"collection_ids"`
	jwt.RegisteredClaims
}

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

func ApiKeyAuthMiddleware(apiKeyUsecase usecase.ApiKeyUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// API keyをヘッダーから取得
		apiKey := c.Request.Header.Get("X-Api-Key")

		// API keyの存在を確認
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-API-Key header is required"})
			c.Abort()
			return
		}

		// API keyの検証
		token, err := apiKeyUsecase.ValidateApiKey(apiKey)
		if err != nil {
			domainErr := &myerrors.DomainError{}
			if errors.As(err, &domainErr) {
				err := controllers.ErrorHandle(domainErr)
				c.JSON(controllers.HttpStatusCodeFromConnectCode(err.Code()), gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}

		// JWT tokenの検証
		claims := &ApiKeyClaims{}
		_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// API keyの検証に成功した場合、userID、projectID、collectionIdsをコンテキストに保存
		c.Set("userID", claims.UserID.String())
		c.Set("projectID", claims.ProjectID)
		c.Set("collectionIds", claims.CollectionIds)

		// API keyが有効な場合、次のハンドラーに進む
		c.Next()
	}
}
