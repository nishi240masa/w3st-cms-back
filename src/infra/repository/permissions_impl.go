package infrastructure

import (
	"context"
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type PermissionRepositoryImpl struct {
	db *gorm.DB
}

func NewPermissionRepositoryImpl(db *gorm.DB) repositories.PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

func (r *PermissionRepositoryImpl) Create(ctx context.Context, permission *models.UserPermission) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Create(permission)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return nil
}

func (r *PermissionRepositoryImpl) FindByID(ctx context.Context, id string) (*models.UserPermission, *myerrors.DomainError) {
	var permission models.UserPermission
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&permission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "権限が見つかりません")
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return &permission, nil
}

func (r *PermissionRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*models.UserPermission, *myerrors.DomainError) {
	var permissions []*models.UserPermission
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&permissions)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return permissions, nil
}

func (r *PermissionRepositoryImpl) FindByUserIDAndResource(ctx context.Context, userID, resource string) ([]*models.UserPermission, *myerrors.DomainError) {
	var permissions []*models.UserPermission
	result := r.db.WithContext(ctx).Where("user_id = ? AND resource = ?", userID, resource).Find(&permissions)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return permissions, nil
}

func (r *PermissionRepositoryImpl) Delete(ctx context.Context, id string) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.UserPermission{})
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	if result.RowsAffected == 0 {
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "権限が見つかりません")
	}
	return nil
}
