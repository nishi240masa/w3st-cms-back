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

func TestMediaUsecase_Upload_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMediaRepo := mockRepositories.NewMockMediaRepository(ctrl)
	uc := usecase.NewMediaUsecase(mockMediaRepo)

	ctx := context.Background()
	userID := uuid.New()
	name := "test.jpg"
	fileType := "image/jpeg"
	path := "/uploads/test.jpg"
	size := int64(1024)

	mockMediaRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	media, err := uc.Upload(ctx, userID, name, fileType, path, size)

	require.NoError(t, err)
	assert.Equal(t, name, media.Name)
	assert.Equal(t, fileType, media.Type)
	assert.Equal(t, path, media.Path)
	assert.Equal(t, size, media.Size)
	assert.Equal(t, userID, media.UserID)
}

func TestMediaUsecase_Upload_FileTooLarge(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMediaRepo := mockRepositories.NewMockMediaRepository(ctrl)
	uc := usecase.NewMediaUsecase(mockMediaRepo)

	ctx := context.Background()
	userID := uuid.New()
	name := "large.jpg"
	fileType := "image/jpeg"
	path := "/uploads/large.jpg"
	size := int64(20 * 1024 * 1024) // 20MB

	_, err := uc.Upload(ctx, userID, name, fileType, path, size)

	require.Error(t, err)
}

func TestMediaUsecase_GetByID_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMediaRepo := mockRepositories.NewMockMediaRepository(ctrl)
	uc := usecase.NewMediaUsecase(mockMediaRepo)

	ctx := context.Background()
	userID := uuid.New()
	id := uuid.New().String()
	expectedMedia := &models.MediaAsset{
		ID:     uuid.MustParse(id),
		Name:   "test.jpg",
		UserID: userID,
	}

	mockMediaRepo.EXPECT().
		FindByID(ctx, id).
		Return(expectedMedia, nil)

	media, err := uc.GetByID(ctx, userID, id)

	require.NoError(t, err)
	assert.Equal(t, expectedMedia, media)
}

func TestMediaUsecase_GetByUserID_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMediaRepo := mockRepositories.NewMockMediaRepository(ctrl)
	uc := usecase.NewMediaUsecase(mockMediaRepo)

	ctx := context.Background()
	userID := uuid.New()
	expectedMedias := []*models.MediaAsset{
		{
			Name:   "test1.jpg",
			UserID: userID,
		},
	}

	mockMediaRepo.EXPECT().
		FindByUserID(ctx, userID.String()).
		Return(expectedMedias, nil)

	medias, err := uc.GetByUserID(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedMedias, medias)
}

func TestMediaUsecase_Delete_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMediaRepo := mockRepositories.NewMockMediaRepository(ctrl)
	uc := usecase.NewMediaUsecase(mockMediaRepo)

	ctx := context.Background()
	userID := uuid.New()
	id := uuid.New().String()
	media := &models.MediaAsset{
		ID:     uuid.MustParse(id),
		UserID: userID,
	}

	mockMediaRepo.EXPECT().
		FindByID(ctx, id).
		Return(media, nil)

	mockMediaRepo.EXPECT().
		Delete(ctx, id).
		Return(nil)

	err := uc.Delete(ctx, userID, id)

	require.NoError(t, err)
}
