package repositories

import (
	"context"

	"w3st/domain/models"
	"w3st/errors"
)

type AuditRepository interface {
	Create(ctx context.Context, log *models.AuditLog) *errors.DomainError
	FindByID(ctx context.Context, id string) (*models.AuditLog, *errors.DomainError)
	FindByUserID(ctx context.Context, userID string) ([]*models.AuditLog, *errors.DomainError)
	FindByAction(ctx context.Context, action string) ([]*models.AuditLog, *errors.DomainError)
	FindAll(ctx context.Context) ([]*models.AuditLog, *errors.DomainError)
}
