package usecase

import (
	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"github.com/google/uuid"
)

type CollectionsUsecase interface {
	Make(newCollection *models.ApiCollection) error
	GetCollectionByUserId(userId uuid.UUID) ([]models.ApiCollection, error)
	GetCollectionsByCollectionId(collectionId int, userId uuid.UUID) (*models.ApiCollection, error)
}

type collectionsUsecase struct {
	collectionsRepo repositories.CollectionsRepository
}

func NewCollectionsUsecase(collectionsRepo repositories.CollectionsRepository) CollectionsUsecase {
	return &collectionsUsecase{
		collectionsRepo: collectionsRepo,
	}
}

func (c *collectionsUsecase) Make(newCollection *models.ApiCollection) error {
	// コレクションを作成する
	err := c.collectionsRepo.CreateCollection(newCollection)
	if err != nil {
		// エラー処理
		return myerrors.WrapDomainError("collectionsUsecase.Make", err)
	}
	return nil
}

func (c *collectionsUsecase) GetCollectionByUserId(userId uuid.UUID) ([]models.ApiCollection, error) {
	collection, err := c.collectionsRepo.GetCollectionByUserId(userId)
	if err != nil {
		return nil, myerrors.WrapDomainError("collectionsUsecase.GetCollectionByUserId", err)
	}
	return collection, nil
}

func (c *collectionsUsecase) GetCollectionsByCollectionId(collectionId int, userId uuid.UUID) (*models.ApiCollection, error) {
	collection, err := c.collectionsRepo.GetCollectionsByCollectionId(collectionId, userId)
	if err != nil {
		return nil, myerrors.WrapDomainError("collectionsUsecase.GetCollectionsByCollectionId", err)
	}
	return collection, nil
}
