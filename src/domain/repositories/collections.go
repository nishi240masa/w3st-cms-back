package repositories

import "w3st/domain/models"

type CollectionsRepository interface {
	CreateCollection(newCollection *models.ApiCollection) error
}
