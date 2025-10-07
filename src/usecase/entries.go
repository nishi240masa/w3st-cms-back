package usecase

import (
	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
)

type EntriesUsecase interface {
	CreateEntry(newEntry *models.Entry, projectId int) error
	GetEntriesByCollectionId(collectionId int, projectId int) ([]models.Entry, error)
	GetEntriesByCollectionIdForSDK(collectionId int, projectId int, collectionIds []int) ([]models.Entry, error)
}

type entriesUsecase struct {
	entriesRepo        repositories.EntriesRepository
	collectionsUsecase CollectionsUsecase
}

func NewEntriesUsecase(entriesRepo repositories.EntriesRepository, collectionsUsecase CollectionsUsecase) EntriesUsecase {
	return &entriesUsecase{
		entriesRepo:        entriesRepo,
		collectionsUsecase: collectionsUsecase,
	}
}

func (e *entriesUsecase) CreateEntry(newEntry *models.Entry, projectId int) error {
	// Check if collection belongs to project
	_, err := e.collectionsUsecase.GetCollectionsByCollectionId(newEntry.CollectionID, projectId)
	if err != nil {
		return myerrors.WrapDomainError("entriesUsecase.CreateEntry", err)
	}

	err = e.entriesRepo.CreateEntry(newEntry)
	if err != nil {
		return myerrors.WrapDomainError("entriesUsecase.CreateEntry", err)
	}
	return nil
}

func (e *entriesUsecase) GetEntriesByCollectionId(collectionId int, projectId int) ([]models.Entry, error) {
	// Check if collection belongs to project
	_, err := e.collectionsUsecase.GetCollectionsByCollectionId(collectionId, projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("entriesUsecase.GetEntriesByCollectionId", err)
	}

	entries, err := e.entriesRepo.GetEntriesByCollectionIdAndProjectId(collectionId, projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("entriesUsecase.GetEntriesByCollectionId", err)
	}
	return entries, nil
}

func (e *entriesUsecase) GetEntriesByCollectionIdForSDK(collectionId int, projectId int, collectionIds []int) ([]models.Entry, error) {
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

	// Check if collection belongs to project
	_, err := e.collectionsUsecase.GetCollectionsByCollectionId(collectionId, projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("entriesUsecase.GetEntriesByCollectionIdForSDK", err)
	}

	entries, err := e.entriesRepo.GetEntriesByCollectionIdAndProjectId(collectionId, projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("entriesUsecase.GetEntriesByCollectionIdForSDK", err)
	}
	return entries, nil
}
