package repositories

import (
	"context"

	"w3st/domain/models"
	"w3st/errors"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *models.UserPermission) *errors.DomainError
	FindByID(ctx context.Context, id string) (*models.UserPermission, *errors.DomainError)
	FindByUserID(ctx context.Context, userID string) ([]*models.UserPermission, *errors.DomainError)
	FindByUserIDAndResource(ctx context.Context, userID, resource string) ([]*models.UserPermission, *errors.DomainError)
	Delete(ctx context.Context, id string) *errors.DomainError
}
