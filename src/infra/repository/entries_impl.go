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

func (r *EntriesRepository) GetEntryByIdAndProjectId(entryId int, projectId int) (*models.Entry, error) {
	var entry models.Entry
	result := r.db.Where("id = ? AND project_id = ?", entryId, projectId).First(&entry)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, result.Error)
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return &entry, nil
}

func (r *EntriesRepository) UpdateEntry(entry *models.Entry) error {
	result := r.db.Save(entry)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return nil
}

func (r *EntriesRepository) DeleteEntry(entryId int, projectId int) error {
	result := r.db.Where("id = ? AND project_id = ?", entryId, projectId).Delete(&models.Entry{})

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	if result.RowsAffected == 0 {
		return myerrors.NewDomainError(myerrors.QueryDataNotFoundError, gorm.ErrRecordNotFound)
	}

	return nil
}
