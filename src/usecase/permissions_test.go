package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"w3st/domain/models"
	mockRepositories "w3st/mock/repositories"
	"w3st/usecase"
)

func TestPermissionUsecase_CheckPermission_HasPermission(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermissionRepo := mockRepositories.NewMockPermissionRepository(ctrl)
	uc := usecase.NewPermissionUsecase(mockPermissionRepo)

	ctx := context.Background()
	userID := uuid.New()
	permission := "read"
	resource := "document"

	permissions := []*models.UserPermission{
		{
			UserID:     userID,
			Permission: permission,
			Resource:   resource,
		},
	}

	mockPermissionRepo.EXPECT().
		FindByUserIDAndResource(ctx, userID.String(), resource).
		Return(permissions, nil)

	hasPermission, err := uc.CheckPermission(ctx, userID, permission, resource)

	require.NoError(t, err)
	assert.True(t, hasPermission)
}

func TestPermissionUsecase_CheckPermission_NoPermission(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermissionRepo := mockRepositories.NewMockPermissionRepository(ctrl)
	uc := usecase.NewPermissionUsecase(mockPermissionRepo)

	ctx := context.Background()
	userID := uuid.New()
	permission := "write"
	resource := "document"

	permissions := []*models.UserPermission{
		{
			UserID:     userID,
			Permission: "read",
			Resource:   resource,
		},
	}

	mockPermissionRepo.EXPECT().
		FindByUserIDAndResource(ctx, userID.String(), resource).
		Return(permissions, nil)

	hasPermission, err := uc.CheckPermission(ctx, userID, permission, resource)

	require.NoError(t, err)
	assert.False(t, hasPermission)
}

func TestPermissionUsecase_GrantPermission_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermissionRepo := mockRepositories.NewMockPermissionRepository(ctrl)
	uc := usecase.NewPermissionUsecase(mockPermissionRepo)

	ctx := context.Background()
	userID := uuid.New()
	permission := "read"
	resource := "document"

	// CheckPermission が false を返す
	mockPermissionRepo.EXPECT().
		FindByUserIDAndResource(ctx, userID.String(), resource).
		Return([]*models.UserPermission{}, nil)

	mockPermissionRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	err := uc.GrantPermission(ctx, userID, permission, resource)

	require.NoError(t, err)
}

func TestPermissionUsecase_GrantPermission_AlreadyExists(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermissionRepo := mockRepositories.NewMockPermissionRepository(ctrl)
	uc := usecase.NewPermissionUsecase(mockPermissionRepo)

	ctx := context.Background()
	userID := uuid.New()
	permission := "read"
	resource := "document"

	permissions := []*models.UserPermission{
		{
			UserID:     userID,
			Permission: permission,
			Resource:   resource,
		},
	}

	mockPermissionRepo.EXPECT().
		FindByUserIDAndResource(ctx, userID.String(), resource).
		Return(permissions, nil)

	err := uc.GrantPermission(ctx, userID, permission, resource)

	require.Error(t, err)
}

func TestPermissionUsecase_RevokePermission_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermissionRepo := mockRepositories.NewMockPermissionRepository(ctrl)
	uc := usecase.NewPermissionUsecase(mockPermissionRepo)

	ctx := context.Background()
	userID := uuid.New()
	permission := "read"
	resource := "document"

	permissions := []*models.UserPermission{
		{
			ID:         uuid.New(),
			UserID:     userID,
			Permission: permission,
			Resource:   resource,
		},
	}

	mockPermissionRepo.EXPECT().
		FindByUserIDAndResource(ctx, userID.String(), resource).
		Return(permissions, nil)

	mockPermissionRepo.EXPECT().
		Delete(ctx, permissions[0].ID.String()).
		Return(nil)

	err := uc.RevokePermission(ctx, userID, permission, resource)

	require.NoError(t, err)
}

func TestPermissionUsecase_GetPermissionsByUser_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermissionRepo := mockRepositories.NewMockPermissionRepository(ctrl)
	uc := usecase.NewPermissionUsecase(mockPermissionRepo)

	ctx := context.Background()
	userID := uuid.New()

	expectedPermissions := []*models.UserPermission{
		{
			UserID:     userID,
			Permission: "read",
			Resource:   "document",
		},
	}

	mockPermissionRepo.EXPECT().
		FindByUserID(ctx, userID.String()).
		Return(expectedPermissions, nil)

	permissions, err := uc.GetPermissionsByUser(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedPermissions, permissions)
}
