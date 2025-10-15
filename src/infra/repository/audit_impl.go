package infrastructure

import (
	"context"
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type AuditRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditRepositoryImpl(db *gorm.DB) repositories.AuditRepository {
	return &AuditRepositoryImpl{db: db}
}

func (r *AuditRepositoryImpl) Create(ctx context.Context, log *models.AuditLog) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Create(log)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return nil
}

func (r *AuditRepositoryImpl) FindByID(ctx context.Context, id string) (*models.AuditLog, *myerrors.DomainError) {
	var log models.AuditLog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&log)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "監査ログが見つかりません")
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return &log, nil
}

func (r *AuditRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*models.AuditLog, *myerrors.DomainError) {
	var logs []*models.AuditLog
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&logs)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return logs, nil
}

func (r *AuditRepositoryImpl) FindByAction(ctx context.Context, action string) ([]*models.AuditLog, *myerrors.DomainError) {
	var logs []*models.AuditLog
	result := r.db.WithContext(ctx).Where("action = ?", action).Find(&logs)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return logs, nil
}

func (r *AuditRepositoryImpl) FindByProjectID(ctx context.Context, projectID int) ([]*models.AuditLog, *myerrors.DomainError) {
	var logs []*models.AuditLog
	result := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&logs)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return logs, nil
}

func (r *AuditRepositoryImpl) FindByProjectIDWithLimit(ctx context.Context, projectID int, limit int, offset int) ([]*models.AuditLog, *myerrors.DomainError) {
	var logs []*models.AuditLog
	result := r.db.WithContext(ctx).Where("project_id = ?", projectID).Limit(limit).Offset(offset).Find(&logs)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return logs, nil
}

func (r *AuditRepositoryImpl) FindAll(ctx context.Context, limit int, offset int) ([]*models.AuditLog, *myerrors.DomainError) {
	var logs []*models.AuditLog
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&logs)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return logs, nil
}
