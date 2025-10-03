package repositories

import (
	"context"

	"w3st/domain/models"
	"w3st/errors"
)

type MediaRepository interface {
	Create(ctx context.Context, media *models.MediaAsset) *errors.DomainError
	FindByID(ctx context.Context, id string) (*models.MediaAsset, *errors.DomainError)
	FindByUserID(ctx context.Context, userID string) ([]*models.MediaAsset, *errors.DomainError)
	Update(ctx context.Context, media *models.MediaAsset) *errors.DomainError
	Delete(ctx context.Context, id string) *errors.DomainError
}
