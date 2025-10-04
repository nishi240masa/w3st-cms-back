package infrastructure

import (
	"w3st/domain/models"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type EntriesRepository struct {
	db *gorm.DB
}

func NewEntriesRepository(db *gorm.DB) *EntriesRepository {
	return &EntriesRepository{
		db: db,
	}
}

func (r *EntriesRepository) CreateEntry(newEntry *models.Entry) error {
	result := r.db.Create(newEntry)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return nil
}

func (r *EntriesRepository) GetEntriesByCollectionIdAndProjectId(collectionId int, projectId int) ([]models.Entry, error) {
	var entries []models.Entry
	result := r.db.Where("collection_id = ? AND project_id = ?", collectionId, projectId).Find(&entries)

	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return entries, nil
}