package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"w3st/domain/models"
	"w3st/domain/repositories"
	"w3st/errors"
)

type ApiKeyClaims struct {
	UserID        uuid.UUID `json:"user_id"`
	ProjectID     int       `json:"project_id"`
	CollectionIds []int     `json:"collection_ids"`
	jwt.RegisteredClaims
}

type JwtUsecase interface {
	GenerateToken(userID uuid.UUID) (models.Token, error)
	ValidateToken(token string) (string, error)
}

type ApiKeyUsecase interface {
	ValidateApiKey(apiKey string) (string, error)
	CreateApiKey(userID string, projectID int, name string, collectionIds []int) (string, error)
}

type jwtAuthUsecase struct {
	secretKey string
}

func NewjwtAuthUsecase() JwtUsecase {
	secret := os.Getenv("SECRET_KEY")
	if len(secret) < 32 {
		panic("SECRET_KEY must be at least 32 bytes long")
	}
	return &jwtAuthUsecase{
		secretKey: secret,
	}
}

// トークンを生成する
func (a *jwtAuthUsecase) GenerateToken(userID uuid.UUID) (models.Token, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, jwtErr := token.SignedString([]byte(a.secretKey))
	if jwtErr != nil {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "トークンの生成に失敗しました")
	}

	return models.Token(signedToken), nil
}

// トークンを検証し、userIDを取得する
func (a *jwtAuthUsecase) ValidateToken(token string) (string, error) {
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "署名方式が不正です")
		}
		return []byte(a.secretKey), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "トークンのパースに失敗しました")
	}

	if (*claims)["exp"] != nil {
		if exp, ok := (*claims)["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "無効なトークンです")
			}
		}
	}

	if (*claims)["sub"] == nil {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "claimsの取得に失敗しました")
	}

	subStr, ok := (*claims)["sub"].(string)
	if !ok {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "subの型が不正です")
	}

	if _, err := uuid.Parse(subStr); err != nil {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "UUIDのパースに失敗しました")
	}

	return subStr, nil
}

type apiKeyUsecase struct {
	repo repositories.ApiKeyRepository
}

func NewApiKeyUsecase(repo repositories.ApiKeyRepository) ApiKeyUsecase {
	return &apiKeyUsecase{repo: repo}
}

func (a *apiKeyUsecase) ValidateApiKey(apiKey string) (string, error) {
	apiKeyModel, err := a.repo.FindByKey(context.Background(), apiKey)
	if err != nil {
		return "", err
	}

	// Generate JWT token with claims
	claims := ApiKeyClaims{
		UserID:        apiKeyModel.UserID,
		ProjectID:     apiKeyModel.ProjectID,
		CollectionIds: apiKeyModel.CollectionIds,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(apiKeyModel.ExpireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, jwtErr := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if jwtErr != nil {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "トークンの生成に失敗しました")
	}

	return signedToken, nil
}

func (a *apiKeyUsecase) CreateApiKey(userID string, projectID int, name string, collectionIds []int) (string, error) {
	// Generate a random API key
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "APIキーの生成に失敗しました")
	}
	apiKey := hex.EncodeToString(bytes)

	// Parse userID to UUID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return "", errors.NewDomainErrorWithMessage(errors.ErrorUnknown, "ユーザーIDのパースに失敗しました")
	}

	// Create API key model
	apiKeyModel := models.ApiKeys{
		UserID:        parsedUserID,
		ProjectID:     projectID,
		Name:          name,
		Key:           apiKey,
		CollectionIds: collectionIds,
		ExpireAt:      time.Now().Add(365 * 24 * time.Hour), // 1 year expiration
		Revoked:       false,
		RateLimit:     1000, // Default rate limit
	}

	// Save to database
	err = a.repo.Create(context.Background(), &apiKeyModel)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}
