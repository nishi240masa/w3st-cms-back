package usecase_test

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"w3st/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupInvalidKey() {
	err := os.Setenv("SECRET_KEY", "short-broken-secret-key-string-for-test!")
	if err != nil {
		panic(err)
	}
}

// func TestGenerateToken_SignedStringFails(t *testing.T) {
// 	t.Parallel()
// 	originalKey := os.Getenv("SECRET_KEY")
// 	defer func() {
// 		os.Setenv("SECRET_KEY", originalKey)
// 	}()

// 	// invalid key type to simulate failure
// 	err := os.Setenv("SECRET_KEY", string([]byte{0xff, 0xfe})) // 署名に失敗しやすくする（強制エラー）
// 	if err != nil {
// 		t.Fatalf("Failed to set SECRET_KEY: %v", err)
// 	}

// 	j := usecase.NewjwtAuthUsecase()

// 	_, err = j.GenerateToken(uuid.New())
// 	assert.Error(t, err)
// }

func TestValidateToken_InvalidSigningMethod(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	// HS256で署名し、algをRS256に変更
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uuid.New().String(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	// トークンを分割し、headerを変更
	parts := strings.Split(signed, ".")
	if len(parts) != 3 {
		t.Fatalf("Invalid token format")
	}
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("Failed to decode header: %v", err)
	}
	var headerMap map[string]interface{}
	if err := json.Unmarshal(header, &headerMap); err != nil {
		t.Fatalf("Failed to unmarshal header: %v", err)
	}
	headerMap["alg"] = "RS256"
	newHeader, err := json.Marshal(headerMap)
	if err != nil {
		t.Fatalf("Failed to marshal header: %v", err)
	}
	newHeaderB64 := base64.RawURLEncoding.EncodeToString(newHeader)
	newToken := newHeaderB64 + "." + parts[1] + "." + parts[2]

	_, err = j.ValidateToken(newToken)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "トークンのパースに失敗しました")
}

func TestValidateToken_ParseError(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	_, err := j.ValidateToken("invalid-token-@@@")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "トークンのパースに失敗しました")
}

func TestValidateToken_InvalidToken(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	// 期限切れトークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uuid.New().String(),
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	_, err := j.ValidateToken(signed)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "無効なトークン")
}

func TestValidateToken_ClaimsMissingSub(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	_, err := j.ValidateToken(signed)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "claims")
}

func TestValidateToken_SubIsNotString(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 123456, // string ではない
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	_, err := j.ValidateToken(signed)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "subの型")
}

func TestValidateToken_SubIsInvalidUUID(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "not-a-uuid",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	_, err := j.ValidateToken(signed)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "UUIDのパース")
}
