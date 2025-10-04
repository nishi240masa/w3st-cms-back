package usecase

import (
	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
)

type CollectionsUsecase interface {
	Make(newCollection *models.ApiCollection) error
	GetCollectionByProjectId(projectId int) ([]models.ApiCollection, error)
	GetCollectionsByCollectionId(collectionId int, projectId int) (*models.ApiCollection, error)
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

func (c *collectionsUsecase) GetCollectionByProjectId(projectId int) ([]models.ApiCollection, error) {
	collection, err := c.collectionsRepo.GetCollectionByProjectId(projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("collectionsUsecase.GetCollectionByProjectId", err)
	}
	return collection, nil
}

func (c *collectionsUsecase) GetCollectionsByCollectionId(collectionId int, projectId int) (*models.ApiCollection, error) {
	collection, err := c.collectionsRepo.GetCollectionsByCollectionId(collectionId, projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("collectionsUsecase.GetCollectionsByCollectionId", err)
	}
	return collection, nil
}
