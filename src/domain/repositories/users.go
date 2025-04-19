package repositories

import (
	"context"

	"w3st/domain/models"
	"w3st/errors"
)

type UserRepository interface {
	Create(ctx context.Context, newUser *models.Users) *errors.DomainError
	FindByEmail(ctx context.Context, email string) (*models.Users, *errors.DomainError)
}
