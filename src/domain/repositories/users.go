package repositories

import (
	"context"

	"w3st/domain/models"
	"w3st/errors"
)

type UserRepository interface {
	Create(ctx context.Context, newUser *models.Users) *errors.DomainError
	FindByEmail(ctx context.Context, email string) (*models.Users, *errors.DomainError)
	FindByID(ctx context.Context, userID string) (*models.Users, *errors.DomainError)
	Update(ctx context.Context, user *models.Users) *errors.DomainError
	GetAllUsers(ctx context.Context) ([]models.Users, *errors.DomainError)
	DeleteUser(ctx context.Context, userID string) *errors.DomainError
}
