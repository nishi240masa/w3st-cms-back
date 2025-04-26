package repositories

import (
	"w3st/domain/models"

	"github.com/google/uuid"
)

type CollectionsRepository interface {
	CreateCollection(newCollection *models.ApiCollection) error
	GetCollectionByUserId(userId uuid.UUID) ([]models.ApiCollection, error)
	GetCollectionsByCollectionId(collectionId string, userId uuid.UUID) (*models.ApiCollection, error)
}
