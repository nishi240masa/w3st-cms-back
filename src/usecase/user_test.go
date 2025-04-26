package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"w3st/domain/models"
	myerrors "w3st/errors"
	mockRepositories "w3st/mock/repositories"
	"w3st/usecase"
)

func TestUserUsecase_Create_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockUserRepo)

	ctx := context.Background()
	newUser := &models.Users{
		ID:       uuid.New(),
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// ユーザーが見つからなかった場合
	mockUserRepo.EXPECT().
		FindByEmail(ctx, "test@example.com").
		Return(nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "not found"))

	mockUserRepo.EXPECT().
		Create(ctx, newUser).
		Return(nil)

	result, err := uc.Create(newUser, ctx)

	require.NoError(t, err)
	assert.Equal(t, newUser, result)
}

func TestUserUsecase_FindByEmail_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockUserRepo)

	email := "success@example.com"
	expectedUser := &models.Users{
		ID:    uuid.New(),
		Email: email,
		Name:  "Successful User",
	}

	mockUserRepo.EXPECT().
		FindByEmail(gomock.Any(), email).
		Return(expectedUser, nil)

	result, err := uc.FindByEmail(email)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, result)
}

func TestUserUsecase_FindById_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockRepositories.NewMockUserRepository(ctrl)
	uc := usecase.NewUserUsecase(mockUserRepo)

	userID := uuid.New()
	expectedUser := &models.Users{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "secret123",
	}

	mockUserRepo.EXPECT().
		FindByID(gomock.Any(), userID.String()).
		Return(expectedUser, nil)

	result, err := uc.FindByID(userID.String())

	require.NoError(t, err)
	assert.Equal(t, expectedUser, result)
}
