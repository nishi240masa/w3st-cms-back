package usecase_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"w3st/usecase"
)

func init() {
	// テスト用のシークレットキーを設定
	err := os.Setenv("SECRET_KEY", "abcdefghijklmnopqrstuvwxyz123456")
	if err != nil {
		return
	}
}

func TestJwtUsecase_GenerateToken_Success(t *testing.T) {
	t.Parallel()
	auth := usecase.NewjwtAuthUsecase()
	userID := uuid.New()

	token, err := auth.GenerateToken(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJwtUsecase_ValidateToken_Success(t *testing.T) {
	t.Parallel()
	auth := usecase.NewjwtAuthUsecase()
	userID := uuid.New()

	token, err := auth.GenerateToken(userID)
	require.NoError(t, err)

	resultID, err := auth.ValidateToken(string(token))
	require.NoError(t, err)
	assert.Equal(t, userID.String(), resultID)
}
