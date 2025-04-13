package usecase_test

import (
	"context"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"w3st/domain/models"
	myerrors "w3st/errors"
	mockRepositories "w3st/mock/repositories"
	mockServices "w3st/mock/services"
	"w3st/usecase"
)

func TestUserUsecase_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	mockAuthService := mockServices.NewMockAuthService(ctrl)
	mockTx := mockRepositories.NewMockTransactionRepository(ctrl)

	uc := usecase.NewUserUsecase(mockUserRepo, mockAuthService, mockTx)

	newUser := &models.Users{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	// モックの期待値設定
	mockTx.
		EXPECT().
		Do(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx) // 実際に渡された関数を実行する
		})

	mockUserRepo.
		EXPECT().
		FindByEmail(gomock.Any(), gomock.Eq("test@example.com")).
		Return(nil, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, "not found"))

	mockUserRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	mockAuthService.
		EXPECT().
		GenerateToken(gomock.Any()).
		Return(models.Token("mocked-token"), nil)

	// テスト実行
	token, err := uc.Create(newUser, context.Background())
	//errの型を確認
	if err != nil {
		log.Fatalf("Create() failed: %+v", err)
	}

	// 検証
	assert.NoError(t, err)
	assert.Equal(t, models.Token("mocked-token"), token)
}
