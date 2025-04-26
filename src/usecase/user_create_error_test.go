package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"w3st/domain/models"
	myerrors "w3st/errors"
	mockRepositories "w3st/mock/repositories"
	"w3st/usecase"
)

func TestUserUsecase_Create_AlreadyExists(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	ctx := context.Background()
	newUser := &models.Users{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	mockRepo.EXPECT().FindByEmail(ctx, "alice@example.com").
		Return(&models.Users{}, nil)

	result, err := uc.Create(newUser, ctx)

	require.Error(t, err)
	assert.Nil(t, result)

	var domainErr *myerrors.DomainError
	require.ErrorAs(t, err, &domainErr)
	assert.Equal(t, myerrors.AlreadyExist, domainErr.ErrType)
}

func TestUserUsecase_Create_FindByEmailFails(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	ctx := context.Background()
	newUser := &models.Users{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "pass123",
	}

	mockRepo.EXPECT().FindByEmail(ctx, "bob@example.com").
		Return(nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryError, "DB接続エラー"))

	result, err := uc.Create(newUser, ctx)

	require.Error(t, err)
	assert.Nil(t, result)

	var domainErr *myerrors.DomainError
	require.ErrorAs(t, err, &domainErr)
	assert.Equal(t, myerrors.QueryError, domainErr.ErrType)
}

func TestUserUsecase_Create_CreateFails(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	ctx := context.Background()
	newUser := &models.Users{
		Name:     "Charlie",
		Email:    "charlie@example.com",
		Password: "pass456",
	}

	mockRepo.EXPECT().FindByEmail(ctx, "charlie@example.com").
		Return(nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "ユーザーが見つかりません"))

	mockRepo.EXPECT().Create(ctx, newUser).
		Return(myerrors.NewDomainErrorWithMessage(myerrors.QueryError, "insert失敗"))

	result, err := uc.Create(newUser, ctx)

	require.Error(t, err)
	assert.Nil(t, result)

	var domainErr *myerrors.DomainError
	require.ErrorAs(t, err, &domainErr)
	assert.Equal(t, myerrors.QueryError, domainErr.ErrType)
}

func TestUserUsecase_FindByEmail_DBError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	email := "notfound@example.com"

	mockRepo.EXPECT().FindByEmail(gomock.Any(), email).
		Return(nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryError, "DB障害"))

	result, err := uc.FindByEmail(email)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestUserUsecase_FindByID_DBError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	userID := "invalid-uuid"

	mockRepo.EXPECT().FindByID(gomock.Any(), userID).
		Return(nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryError, "DB障害"))

	result, err := uc.FindByID(userID)

	require.Error(t, err)
	assert.Nil(t, result)
}
