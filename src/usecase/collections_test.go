package usecase_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"w3st/domain/models"
	mockRepositories "w3st/mock/repositories"
	"w3st/usecase"
)

func TestCollectionsUsecase_Make_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewCollectionsUsecase(mockCollectionsRepo)

	newCollection := &models.ApiCollection{
		Name:   "Test Collection",
		UserID: uuid.New(),
	}

	mockCollectionsRepo.EXPECT().
		CreateCollection(newCollection).
		Return(nil)

	err := uc.Make(newCollection)

	require.NoError(t, err)
}

func TestCollectionsUsecase_Make_Failure(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewCollectionsUsecase(mockCollectionsRepo)

	newCollection := &models.ApiCollection{
		Name:   "Test Collection",
		UserID: uuid.New(),
	}

	mockCollectionsRepo.EXPECT().
		CreateCollection(newCollection).
		Return(errors.New("test error"))

	err := uc.Make(newCollection)

	require.Error(t, err)
}

func TestCollectionsUsecase_GetCollectionByUserId_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewCollectionsUsecase(mockCollectionsRepo)

	userID := uuid.New()
	expectedCollections := []models.ApiCollection{
		{
			Name:   "Collection 1",
			UserID: userID,
		},
	}

	mockCollectionsRepo.EXPECT().
		GetCollectionByUserId(userID).
		Return(expectedCollections, nil)

	collections, err := uc.GetCollectionByUserId(userID)

	require.NoError(t, err)
	assert.Equal(t, expectedCollections, collections)
}

func TestCollectionsUsecase_GetCollectionsByCollectionId_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollectionsRepo := mockRepositories.NewMockCollectionsRepository(ctrl)
	uc := usecase.NewCollectionsUsecase(mockCollectionsRepo)

	collectionID := 1
	userID := uuid.New()
	expectedCollection := &models.ApiCollection{
		ID:     collectionID,
		Name:   "Test Collection",
		UserID: userID,
	}

	mockCollectionsRepo.EXPECT().
		GetCollectionsByCollectionId(collectionID, userID).
		Return(expectedCollection, nil)

	collection, err := uc.GetCollectionsByCollectionId(collectionID, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedCollection, collection)
}
