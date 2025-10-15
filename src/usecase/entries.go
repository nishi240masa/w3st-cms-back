package usecase

import (
	"encoding/json"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
)

type EntriesUsecase interface {
	CreateEntry(newEntry *models.Entry, projectId int) error
	GetEntriesByCollectionId(collectionId int, projectId int) ([]models.Entry, error)
	GetEntriesByCollectionIdForSDK(collectionId int, projectId int, collectionIds []int) ([]models.Entry, error)
	UpdateEntry(entryId int, data map[string]interface{}, projectId int) error
	DeleteEntry(entryId int, projectId int) error
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

func (e *entriesUsecase) UpdateEntry(entryId int, data map[string]interface{}, projectId int) error {
	// Check if entry exists and belongs to project
	entry, err := e.entriesRepo.GetEntryByIdAndProjectId(entryId, projectId)
	if err != nil {
		return myerrors.WrapDomainError("entriesUsecase.UpdateEntry", err)
	}

	// Update entry data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return myerrors.NewDomainError(myerrors.QueryError, err)
	}

	entry.Data = string(dataBytes)

	err = e.entriesRepo.UpdateEntry(entry)
	if err != nil {
		return myerrors.WrapDomainError("entriesUsecase.UpdateEntry", err)
	}
	return nil
}

func (e *entriesUsecase) DeleteEntry(entryId int, projectId int) error {
	// Check if entry exists and belongs to project
	_, err := e.entriesRepo.GetEntryByIdAndProjectId(entryId, projectId)
	if err != nil {
		return myerrors.WrapDomainError("entriesUsecase.DeleteEntry", err)
	}

	err = e.entriesRepo.DeleteEntry(entryId, projectId)
	if err != nil {
		return myerrors.WrapDomainError("entriesUsecase.DeleteEntry", err)
	}
	return nil
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
