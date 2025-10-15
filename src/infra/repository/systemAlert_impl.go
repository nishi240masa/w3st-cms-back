package infrastructure

import (
	"context"
	"errors"

	"w3st/domain/models"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type SystemAlertRepository struct {
	db *gorm.DB
}

func NewSystemAlertRepository(db *gorm.DB) *SystemAlertRepository {
	return &SystemAlertRepository{
		db: db,
	}
}

func (r *SystemAlertRepository) Create(ctx context.Context, alert *models.SystemAlert) error {
	result := r.db.WithContext(ctx).Create(alert)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return nil
}

func (r *SystemAlertRepository) FindByID(ctx context.Context, id int) (*models.SystemAlert, error) {
	var alert models.SystemAlert
	result := r.db.WithContext(ctx).First(&alert, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, result.Error)
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return &alert, nil
}

func (r *SystemAlertRepository) FindActiveByProjectID(ctx context.Context, projectID int) ([]models.SystemAlert, error) {
	var alerts []models.SystemAlert
	result := r.db.WithContext(ctx).Where("project_id = ? AND is_active = ?", projectID, true).Order("created_at DESC").Find(&alerts)

	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return alerts, nil
}

func (r *SystemAlertRepository) FindAllByProjectID(ctx context.Context, projectID int, limit int, offset int) ([]models.SystemAlert, error) {
	var alerts []models.SystemAlert
	result := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&alerts)

	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return alerts, nil
}

func (r *SystemAlertRepository) Update(ctx context.Context, alert *models.SystemAlert) error {
	result := r.db.WithContext(ctx).Save(alert)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return nil
}

func (r *SystemAlertRepository) MarkAsRead(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Model(&models.SystemAlert{}).Where("id = ?", id).Update("is_read", true)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	if result.RowsAffected == 0 {
		return myerrors.NewDomainError(myerrors.QueryDataNotFoundError, gorm.ErrRecordNotFound)
	}

	return nil
}

func (r *SystemAlertRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.SystemAlert{}, id)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	if result.RowsAffected == 0 {
		return myerrors.NewDomainError(myerrors.QueryDataNotFoundError, gorm.ErrRecordNotFound)
	}

	return nil
}

func (r *SystemAlertRepository) CountActiveByProjectID(ctx context.Context, projectID int) (int, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.SystemAlert{}).Where("project_id = ? AND is_active = ?", projectID, true).Count(&count)

	if result.Error != nil {
		return 0, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return int(count), nil
}
