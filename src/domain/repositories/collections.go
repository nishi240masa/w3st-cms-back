package repositories

import (
	"w3st/domain/models"
)

type CollectionsRepository interface {
	CreateCollection(newCollection *models.ApiCollection) error
	GetCollectionByProjectId(projectId int) ([]models.ApiCollection, error)
	GetCollectionsByCollectionId(collectionId int, projectId int) (*models.ApiCollection, error)
}
