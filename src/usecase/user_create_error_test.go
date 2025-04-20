package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"w3st/domain/models"
	"w3st/usecase"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	myerrors "w3st/errors"
	mockRepositories "w3st/mock/repositories"
	mockServices "w3st/mock/services"
)

func TestUserUsecase_Create_AlreadyExists(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	mockAuthService := mockServices.NewMockAuthService(ctrl)
	mockTx := mockRepositories.NewMockTransactionRepository(ctrl)

	uc := usecase.NewUserUsecase(mockUserRepo, mockAuthService, mockTx)

	newUser := &models.Users{
		Name:     "Existing User",
		Email:    "exist@example.com",
		Password: "pass123",
	}

	mockTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	mockUserRepo.EXPECT().FindByEmail(gomock.Any(), "exist@example.com").
		Return(&models.Users{}, nil)

	token, err := uc.Create(newUser, context.Background())

	require.Error(t, err)
	assert.Equal(t, models.Token(""), token)
}

func TestUserUsecase_Create_TokenGenerationFails(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	mockAuthService := mockServices.NewMockAuthService(ctrl)
	mockTx := mockRepositories.NewMockTransactionRepository(ctrl)

	uc := usecase.NewUserUsecase(mockUserRepo, mockAuthService, mockTx)

	newUser := &models.Users{
		Name:     "Token Error",
		Email:    "tokenfail@example.com",
		Password: "pass123",
	}

	mockTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	mockUserRepo.EXPECT().FindByEmail(gomock.Any(), "tokenfail@example.com").
		Return(nil, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, "not found"))

	mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		Return(nil)

	mockAuthService.EXPECT().GenerateToken(gomock.Any()).
		Return(models.Token(""), myerrors.NewDomainError(myerrors.RepositoryError, "token生成失敗"))
	token, err := uc.Create(newUser, context.Background())

	require.Error(t, err)
	assert.Equal(t, models.Token(""), token)
}

func TestUserUsecase_FindByEmail_RepoFails(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	mockAuthService := mockServices.NewMockAuthService(ctrl)
	mockTx := mockRepositories.NewMockTransactionRepository(ctrl)

	uc := usecase.NewUserUsecase(mockUserRepo, mockAuthService, mockTx)

	email := "notfound@example.com"

	mockUserRepo.EXPECT().
		FindByEmail(gomock.Any(), email).
		Return(nil, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, "not found"))

	token, err := uc.FindByEmail(email)

	require.Error(t, err)
	assert.Equal(t, models.Token(""), token)
}

func TestUserUsecase_FindByEmail_QueryError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)

	uc := usecase.NewUserUsecase(mockUserRepo, nil, nil)

	email := "queryfail@example.com"

	mockUserRepo.EXPECT().
		FindByEmail(gomock.Any(), email).
		Return(nil, myerrors.NewDomainError(myerrors.QueryError, "DB error"))

	token, err := uc.FindByEmail(email)
	require.Error(t, err)
	assert.Equal(t, models.Token(""), token)
}

func TestUserUsecase_FindByEmail_TokenGenerationFails(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	mockAuthService := mockServices.NewMockAuthService(ctrl)
	mockTx := mockRepositories.NewMockTransactionRepository(ctrl)

	uc := usecase.NewUserUsecase(mockUserRepo, mockAuthService, mockTx)

	email := "failtoken@example.com"
	userID := uuid.New()
	mockUser := &models.Users{
		ID: userID,
	}

	mockUserRepo.EXPECT().
		FindByEmail(gomock.Any(), email).
		Return(mockUser, nil)

	mockAuthService.EXPECT().
		GenerateToken(userID).
		Return(models.Token(""), myerrors.NewDomainError(myerrors.RepositoryError, "token生成失敗"))

	token, err := uc.FindByEmail(email)

	require.Error(t, err)
	assert.Equal(t, models.Token(""), token)
}
