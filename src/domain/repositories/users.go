package repositories

import (
	"context"
	"w3st/domain/models"
	"w3st/errors"
)

type UserRepository interface {
	Create(ctx context.Context,newUser *models.User)  *errors.DomainError
	FindByEmail(ctx context.Context,email string) (*models.User,  *errors.DomainError)
}