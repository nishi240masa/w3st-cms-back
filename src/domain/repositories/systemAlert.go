package repositories

import (
	"context"
	"w3st/domain/models"
)

type SystemAlertRepository interface {
	Create(ctx context.Context, alert *models.SystemAlert) error
	FindByID(ctx context.Context, id int) (*models.SystemAlert, error)
	FindActiveByProjectID(ctx context.Context, projectID int) ([]models.SystemAlert, error)
	FindAllByProjectID(ctx context.Context, projectID int, limit int, offset int) ([]models.SystemAlert, error)
	Update(ctx context.Context, alert *models.SystemAlert) error
	MarkAsRead(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
	CountActiveByProjectID(ctx context.Context, projectID int) (int, error)
}