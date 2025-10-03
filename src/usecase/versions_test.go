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

func TestVersionUsecase_CreateVersion_FirstVersion(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mockRepositories.NewMockVersionRepository(ctrl)
	uc := usecase.NewVersionUsecase(mockVersionRepo)

	ctx := context.Background()
	userID := uuid.New()
	contentID := uuid.New()
	data := map[string]string{"key": "value"}

	mockVersionRepo.EXPECT().
		FindLatestByContentID(ctx, contentID.String()).
		Return(nil, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, nil))

	mockVersionRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	version, err := uc.CreateVersion(ctx, userID, contentID, data)

	require.NoError(t, err)
	assert.Equal(t, contentID, version.ContentID)
	assert.Equal(t, 1, version.Version)
	assert.Equal(t, userID, version.UserID)
}

func TestVersionUsecase_GetVersionsByContentID_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mockRepositories.NewMockVersionRepository(ctrl)
	uc := usecase.NewVersionUsecase(mockVersionRepo)

	ctx := context.Background()
	userID := uuid.New()
	contentID := uuid.New()

	expectedVersions := []*models.ContentVersion{
		{
			ContentID: contentID,
			Version:   1,
			UserID:    userID,
		},
	}

	mockVersionRepo.EXPECT().
		FindByContentID(ctx, contentID.String()).
		Return(expectedVersions, nil)

	versions, err := uc.GetVersionsByContentID(ctx, userID, contentID)

	require.NoError(t, err)
	assert.Equal(t, expectedVersions, versions)
}

func TestVersionUsecase_GetLatestVersion_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mockRepositories.NewMockVersionRepository(ctrl)
	uc := usecase.NewVersionUsecase(mockVersionRepo)

	ctx := context.Background()
	userID := uuid.New()
	contentID := uuid.New()

	expectedVersion := &models.ContentVersion{
		ContentID: contentID,
		Version:   2,
		UserID:    userID,
	}

	mockVersionRepo.EXPECT().
		FindLatestByContentID(ctx, contentID.String()).
		Return(expectedVersion, nil)

	version, err := uc.GetLatestVersion(ctx, userID, contentID)

	require.NoError(t, err)
	assert.Equal(t, expectedVersion, version)
}

func TestVersionUsecase_RestoreVersion_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mockRepositories.NewMockVersionRepository(ctrl)
	uc := usecase.NewVersionUsecase(mockVersionRepo)

	ctx := context.Background()
	userID := uuid.New()
	contentID := uuid.New()
	versionID := uuid.New()

	expectedVersion := &models.ContentVersion{
		ID:        versionID,
		ContentID: contentID,
		Version:   1,
		UserID:    userID,
	}

	mockVersionRepo.EXPECT().
		FindByID(ctx, versionID.String()).
		Return(expectedVersion, nil)

	version, err := uc.RestoreVersion(ctx, userID, contentID, versionID)

	require.NoError(t, err)
	assert.Equal(t, expectedVersion, version)
}
