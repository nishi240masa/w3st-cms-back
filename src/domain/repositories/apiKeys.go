package repositories

import (
	"context"

	"w3st/domain/models"
	"w3st/errors"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey *models.ApiKeys) *errors.DomainError
	FindByKey(ctx context.Context, key string) (*models.ApiKeys, *errors.DomainError)
	FindByUserID(ctx context.Context, userID string) ([]*models.ApiKeys, *errors.DomainError)
	Update(ctx context.Context, apiKey *models.ApiKeys) *errors.DomainError
	Delete(ctx context.Context, id int) *errors.DomainError
}