package usecase_test

import (
	"context"
	"errors"
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

func TestAuditUsecase_LogAction_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuditRepo := mockRepositories.NewMockAuditRepository(ctrl)
	uc := usecase.NewAuditUsecase(mockAuditRepo)

	ctx := context.Background()
	userID := uuid.New()
	action := "create"
	resource := "user"
	details := "Created new user"

	mockAuditRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	err := uc.LogAction(ctx, userID, action, resource, details)

	require.NoError(t, err)
}

func TestAuditUsecase_LogAction_Failure(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuditRepo := mockRepositories.NewMockAuditRepository(ctrl)
	uc := usecase.NewAuditUsecase(mockAuditRepo)

	ctx := context.Background()
	userID := uuid.New()
	action := "create"
	resource := "user"
	details := "Created new user"

	mockAuditRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(myerrors.NewDomainError(myerrors.RepositoryError, errors.New("test error")))

	err := uc.LogAction(ctx, userID, action, resource, details)

	require.Error(t, err)
}

func TestAuditUsecase_GetLogsByUser_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuditRepo := mockRepositories.NewMockAuditRepository(ctrl)
	uc := usecase.NewAuditUsecase(mockAuditRepo)

	ctx := context.Background()
	userID := uuid.New()
	expectedLogs := []*models.AuditLog{
		{
			UserID:   userID,
			Action:   "login",
			Resource: "auth",
			Details:  "User logged in",
		},
	}

	mockAuditRepo.EXPECT().
		FindByUserID(ctx, userID.String()).
		Return(expectedLogs, nil)

	logs, err := uc.GetLogsByUser(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedLogs, logs)
}

func TestAuditUsecase_GetLogsByAction_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuditRepo := mockRepositories.NewMockAuditRepository(ctrl)
	uc := usecase.NewAuditUsecase(mockAuditRepo)

	ctx := context.Background()
	action := "create"
	expectedLogs := []*models.AuditLog{
		{
			Action:   action,
			Resource: "user",
			Details:  "Created new user",
		},
	}

	mockAuditRepo.EXPECT().
		FindByAction(ctx, action).
		Return(expectedLogs, nil)

	logs, err := uc.GetLogsByAction(ctx, action)

	require.NoError(t, err)
	assert.Equal(t, expectedLogs, logs)
}

func TestAuditUsecase_GetAllLogs_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuditRepo := mockRepositories.NewMockAuditRepository(ctrl)
	uc := usecase.NewAuditUsecase(mockAuditRepo)

	ctx := context.Background()
	expectedLogs := []*models.AuditLog{
		{
			Action:   "login",
			Resource: "auth",
			Details:  "User logged in",
		},
	}

	mockAuditRepo.EXPECT().
		FindAll(ctx).
		Return(expectedLogs, nil)

	logs, err := uc.GetAllLogs(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedLogs, logs)
}