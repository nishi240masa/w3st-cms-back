package repositories

import (
	"w3st/domain/models"
)

type EntriesRepository interface {
	CreateEntry(newEntry *models.Entry) error
	GetEntriesByCollectionIdAndProjectId(collectionId int, projectId int) ([]models.Entry, error)
	GetEntryByIdAndProjectId(entryId int, projectId int) (*models.Entry, error)
	UpdateEntry(entry *models.Entry) error
	DeleteEntry(entryId int, projectId int) error
}
