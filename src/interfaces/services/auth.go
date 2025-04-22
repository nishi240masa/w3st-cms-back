package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"w3st/domain/models"
	"w3st/errors"
)

type AuthService interface {
	GenerateToken(userID uuid.UUID) (models.Token, *errors.DomainError)
	ValidateToken(token string) (string, *errors.DomainError)
}

type authService struct {
	secretKey string
}

func NewAuthService() AuthService {
	secret := os.Getenv("SECRET_KEY")
	if len(secret) < 32 {
		panic("SECRET_KEY must be at least 32 bytes long")
	}
	return &authService{
		secretKey: secret,
	}
}

// トークンを生成する
func (a *authService) GenerateToken(userID uuid.UUID) (models.Token, *errors.DomainError) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", errors.NewDomainError(errors.ErrorUnknown, "トークンの生成に失敗しました")
	}

	return models.Token(signedToken), nil
}

// トークンを検証し、userIDを取得する
func (a *authService) ValidateToken(token string) (string, *errors.DomainError) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewDomainError(errors.ErrorUnknown, "署名方式が不正です")
		}
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return "", errors.NewDomainError(errors.ErrorUnknown, "トークンのパースに失敗しました")
	}

	if !parsedToken.Valid {
		return "", errors.NewDomainError(errors.ErrorUnknown, "無効なトークンです")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return "", errors.NewDomainError(errors.ErrorUnknown, "claimsの取得に失敗しました")
	}

	subStr, ok := claims["sub"].(string)
	if !ok {
		return "", errors.NewDomainError(errors.ErrorUnknown, "subの型が不正です")
	}

	if _, err := uuid.Parse(subStr); err != nil {
		return "", errors.NewDomainError(errors.ErrorUnknown, "UUIDのパースに失敗しました")
	}

	return subStr, nil
}
