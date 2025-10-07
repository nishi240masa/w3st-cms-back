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
	GetCollectionByProjectIdForSDK(projectId int, collectionIds []int) ([]models.ApiCollection, error)
	GetCollectionsByCollectionIdForSDK(collectionId int, projectId int, collectionIds []int) (*models.ApiCollection, error)
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

func (c *collectionsUsecase) GetCollectionByProjectIdForSDK(projectId int, collectionIds []int) ([]models.ApiCollection, error) {
	// Get all collections for the project
	allCollections, err := c.collectionsRepo.GetCollectionByProjectId(projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("collectionsUsecase.GetCollectionByProjectIdForSDK", err)
	}

	// Filter collections based on allowed collectionIds
	var filteredCollections []models.ApiCollection
	for _, collection := range allCollections {
		for _, allowedId := range collectionIds {
			if collection.ID == allowedId {
				filteredCollections = append(filteredCollections, collection)
				break
			}
		}
	}

	return filteredCollections, nil
}

func (c *collectionsUsecase) GetCollectionsByCollectionIdForSDK(collectionId int, projectId int, collectionIds []int) (*models.ApiCollection, error) {
	// Check if collectionId is in allowed collectionIds
	allowed := false
	for _, id := range collectionIds {
		if id == collectionId {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "Collection not accessible with this API key")
	}

	collection, err := c.collectionsRepo.GetCollectionsByCollectionId(collectionId, projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("collectionsUsecase.GetCollectionsByCollectionIdForSDK", err)
	}
	return collection, nil
}
