package usecase_test

import (
	"os"
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

func TestGenerateToken_SignedStringFails(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	// invalid key type to simulate failure
	originalKey := os.Getenv("SECRET_KEY")
	if os.Getenv("SECRET_KEY") == originalKey {
		t.Fatalf("SECRET_KEY should be restored after test")
	}

	err := os.Setenv("SECRET_KEY", string([]byte{0xff, 0xfe})) // 署名に失敗しやすくする（強制エラー）
	if err != nil {
		t.Fatalf("Failed to set SECRET_KEY: %v", err)
	}

	_, err = j.GenerateToken(uuid.New())
	assert.Error(t, err)
}

func TestValidateToken_InvalidSigningMethod(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	// RS256で不正に署名
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": uuid.New().String(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte("dummy"))

	_, err := j.ValidateToken(signed)
	require.NoError(t, err)
	assert.Contains(t, err.Error(), "署名方式")
}

func TestValidateToken_ParseError(t *testing.T) {
	t.Parallel()
	setupInvalidKey()
	j := usecase.NewjwtAuthUsecase()

	_, err := j.ValidateToken("invalid-token-@@@")
	require.NoError(t, err)
	assert.Contains(t, err.Error(), "パース")
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
	require.NoError(t, err)
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
	require.NoError(t, err)
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
	require.NoError(t, err)
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
	require.NoError(t, err)
	assert.Contains(t, err.Error(), "UUIDのパース")
}
