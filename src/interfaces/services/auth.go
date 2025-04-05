package services

import (
	"os"
	"time"
	"w3st/domain/models"
	"w3st/errors"

	"github.com/golang-jwt/jwt/v5"
)


type AuthService interface {
    GenerateToken(userID string) (models.Token,  *errors.DomainError)
    ValidateToken(token string) (string, error)
}

type authService struct {
	secretKey string
}

func NewAuthService() AuthService {
	return &authService{
		secretKey: os.Getenv("SECRET_KEY"),
	}
}

// tokenの生成
func (a *authService) GenerateToken(userID string) (models.Token,  *errors.DomainError) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// トークンの署名
	signedToken, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "",  errors.NewDomainError(errors.ErrorUnknown, "トークンの生成に失敗しました")
	}
	return models.Token(signedToken), nil
}

// tokenの検証
func (a *authService) ValidateToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			// 署名方法がHS256でない場合
			return nil, errors.NewDomainError(errors.ErrorUnknown, "invalid signing method")
		}
		return []byte(a.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		// トークンが無効な場合
		return "", errors.NewDomainError(errors.ErrorUnknown, "invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		// claimsが無効な場合
		return "", errors.NewDomainError(errors.ErrorUnknown, "invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		// user_idがstring型でない場合
		return "", errors.NewDomainError(errors.ErrorUnknown, "invalid user_id")
	}

	return userID, nil
}
