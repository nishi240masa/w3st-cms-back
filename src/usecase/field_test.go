package usecase_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"w3st/domain/models"
	mockRepositories "w3st/mock/repositories"
	"w3st/usecase"
)

func TestFieldUsecase_Create_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFieldRepo := mockRepositories.NewMockFieldRepository(ctrl)
	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewFieldUsecase(mockFieldRepo, mockCollectionsRepo)

	userID := uuid.New()
	newField := &models.FieldData{
		CollectionID: 1,
		ViewName:     "Test Field",
	}

	mockCollectionsRepo.EXPECT().
		GetCollectionsByCollectionId(1, userID).
		Return(&models.ApiCollection{}, nil)

	mockFieldRepo.EXPECT().
		CreateField(newField).
		Return(nil)

	err := uc.Create(userID, newField)

	require.NoError(t, err)
}

func TestFieldUsecase_Create_CollectionNotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFieldRepo := mockRepositories.NewMockFieldRepository(ctrl)
	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewFieldUsecase(mockFieldRepo, mockCollectionsRepo)

	userID := uuid.New()
	newField := &models.FieldData{
		CollectionID: 1,
		ViewName:     "Test Field",
	}

	mockCollectionsRepo.EXPECT().
		GetCollectionsByCollectionId(1, userID).
		Return(nil, errors.New("collection not found"))

	err := uc.Create(userID, newField)

	require.Error(t, err)
}

func TestFieldUsecase_Update_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFieldRepo := mockRepositories.NewMockFieldRepository(ctrl)
	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewFieldUsecase(mockFieldRepo, mockCollectionsRepo)

	userID := uuid.New()
	newField := &models.FieldData{
		CollectionID: 1,
		ViewName:     "Updated Field",
	}

	mockCollectionsRepo.EXPECT().
		GetCollectionsByCollectionId(1, userID).
		Return(&models.ApiCollection{}, nil)

	mockFieldRepo.EXPECT().
		UpdateField(newField).
		Return(nil)

	err := uc.Update(userID, newField)

	require.NoError(t, err)
}

func TestFieldUsecase_Delete_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFieldRepo := mockRepositories.NewMockFieldRepository(ctrl)
	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewFieldUsecase(mockFieldRepo, mockCollectionsRepo)

	userID := uuid.New().String()
	fieldID := uuid.New().String()

	mockFieldRepo.EXPECT().
		DeleteFieldById(uuid.MustParse(userID), uuid.MustParse(fieldID)).
		Return(nil)

	err := uc.Delete(userID, fieldID)

	require.NoError(t, err)
}

func TestFieldUsecase_Delete_InvalidUUID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFieldRepo := mockRepositories.NewMockFieldRepository(ctrl)
	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewFieldUsecase(mockFieldRepo, mockCollectionsRepo)

	userID := "testInvalidUUID"
	fieldID := uuid.New().String()

	err := uc.Delete(userID, fieldID)

	require.Error(t, err)
}
