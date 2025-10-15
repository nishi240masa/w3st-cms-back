package middlewares

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	myerrors "w3st/errors"
	"w3st/interfaces/controllers"
	"w3st/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ApiKeyClaims struct {
	UserID        uuid.UUID `json:"user_id"`
	ProjectID     int       `json:"project_id"`
	CollectionIds []int     `json:"collection_ids"`
	jwt.RegisteredClaims
}

type Auth0Claims struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.RegisteredClaims
}

type Auth0JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	X5c []string `json:"x5c"`
}

type Auth0JWKSet struct {
	Keys []Auth0JWK `json:"keys"`
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

// ProjectRateLimitMiddleware プロジェクト単位のレート制限ミドルウェア
func ProjectRateLimitMiddleware(projectUsecase usecase.ProjectUsecase, systemAlertUsecase usecase.SystemAlertUsecase) gin.HandlerFunc {
	// インメモリストレージ（本番ではRedisなどの外部ストレージを使用）
	requestCounts := make(map[string]int)
	lastReset := make(map[string]time.Time)

	return func(c *gin.Context) {
		projectID := c.GetInt("projectID")
		if projectID == 0 {
			c.Next()
			return
		}

		// プロジェクトのレート制限を取得
		rateLimit, err := projectUsecase.GetRateLimitByProjectID(c.Request.Context(), projectID)
		if err != nil {
			// レート制限が取得できない場合はデフォルト値を使用
			rateLimit = 1000
		}

		// 現在の時間キー（1時間単位）
		timeKey := time.Now().Format("2006-01-02-15")
		mapKey := fmt.Sprintf("project_%d_%s", projectID, timeKey)

		// 前回のリセット時間をチェック（1時間ごとにリセット）
		if lastTime, exists := lastReset[mapKey]; exists {
			if time.Since(lastTime) >= time.Hour {
				// 1時間経過したらカウントをリセット
				delete(requestCounts, mapKey)
				delete(lastReset, mapKey)
			}
		}

		// リクエストカウントを取得・インクリメント
		count := requestCounts[mapKey]
		count++
		requestCounts[mapKey] = count
		lastReset[mapKey] = time.Now()

		// レート制限チェック
		if count > rateLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded for this project",
				"limit": rateLimit,
				"reset": time.Now().Add(time.Hour).Format(time.RFC3339),
			})
			c.Abort()
			return
		}

		// アラートチェック（100リクエストごとに）
			if count%100 == 0 {
				err := systemAlertUsecase.CheckAndCreateApiLimitAlert(c.Request.Context(), projectID, count, rateLimit)
				if err != nil {
					// アラート作成失敗はログ出力のみ（リクエストは継続）
				}
			}
	
			c.Next()
		}
	}
	
	// Auth0AuthMiddleware Auth0トークン検証ミドルウェア
	func Auth0AuthMiddleware() gin.HandlerFunc {
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
	
			// Auth0トークンの検証
			claims, err := validateAuth0Token(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Auth0 token"})
				c.Abort()
				return
			}
	
			// 検証成功の場合、ユーザー情報をコンテキストに保存
			c.Set("userID", claims.Sub)
			c.Set("userEmail", claims.Email)
			c.Set("userName", claims.Name)
	
			c.Next()
		}
	}
	
	// validateAuth0Token Auth0トークンを検証
	func validateAuth0Token(tokenString string) (*Auth0Claims, error) {
		// Auth0ドメインを取得
		auth0Domain := os.Getenv("AUTH0_DOMAIN")
		if auth0Domain == "" {
			return nil, errors.New("AUTH0_DOMAIN environment variable is not set")
		}
	
		// JWK Setを取得
		jwkSet, err := getAuth0JWKSet(auth0Domain)
		if err != nil {
			return nil, err
		}
	
		// トークンをパース（署名検証なし）
		token, err := jwt.ParseWithClaims(tokenString, &Auth0Claims{}, nil)
		if err != nil {
			return nil, err
		}
	
		if _, ok := token.Claims.(*Auth0Claims); !ok {
			return nil, errors.New("invalid token claims")
		}
	
		// kidから適切な公開鍵を取得
		kid := token.Header["kid"].(string)
		var publicKey *rsa.PublicKey
		for _, jwk := range jwkSet.Keys {
			if jwk.Kid == kid {
				publicKey, err = jwkToPublicKey(jwk)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	
		if publicKey == nil {
			return nil, errors.New("unable to find appropriate key")
		}
	
		// 署名検証付きで再度パース
		token, err = jwt.ParseWithClaims(tokenString, &Auth0Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
	
		if err != nil {
			return nil, err
		}
	
		if claims, ok := token.Claims.(*Auth0Claims); ok && token.Valid {
			return claims, nil
		}
	
		return nil, errors.New("invalid token")
	}
	
	// getAuth0JWKSet Auth0からJWK Setを取得
	func getAuth0JWKSet(domain string) (*Auth0JWKSet, error) {
		resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", domain))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	
		var jwkSet Auth0JWKSet
		if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
			return nil, err
		}
	
		return &jwkSet, nil
	}
	
	// jwkToPublicKey JWKからRSA公開鍵を生成
	func jwkToPublicKey(jwk Auth0JWK) (*rsa.PublicKey, error) {
		if jwk.Kty != "RSA" {
			return nil, errors.New("unsupported key type")
		}
	
		// Base64URLデコード
		nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			return nil, err
		}
	
		eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			return nil, err
		}
	
		// RSA公開鍵を構築
		publicKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
	
		return publicKey, nil
	}
