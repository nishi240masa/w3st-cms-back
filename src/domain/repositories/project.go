package repositories

import (
	"context"

	"w3st/domain/models"
)

type ProjectRepository interface {
	FindByID(ctx context.Context, id int) (*models.Project, error)
	FindAll(ctx context.Context) ([]models.Project, error)
	Create(ctx context.Context, project *models.Project) (*models.Project, error)
	Update(ctx context.Context, project *models.Project) error
}
