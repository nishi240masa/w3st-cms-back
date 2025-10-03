package repositories

import (
	"context"

	"w3st/domain/models"
	"w3st/errors"
)

type VersionRepository interface {
	Create(ctx context.Context, version *models.ContentVersion) *errors.DomainError
	FindByID(ctx context.Context, id string) (*models.ContentVersion, *errors.DomainError)
	FindByContentID(ctx context.Context, contentID string) ([]*models.ContentVersion, *errors.DomainError)
	FindLatestByContentID(ctx context.Context, contentID string) (*models.ContentVersion, *errors.DomainError)
	Delete(ctx context.Context, id string) *errors.DomainError
}
