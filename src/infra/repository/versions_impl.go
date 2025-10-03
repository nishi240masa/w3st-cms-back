package infrastructure

import (
	"context"
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type VersionRepositoryImpl struct {
	db *gorm.DB
}

func NewVersionRepositoryImpl(db *gorm.DB) repositories.VersionRepository {
	return &VersionRepositoryImpl{db: db}
}

func (r *VersionRepositoryImpl) Create(ctx context.Context, version *models.ContentVersion) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Create(version)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return nil
}

func (r *VersionRepositoryImpl) FindByID(ctx context.Context, id string) (*models.ContentVersion, *myerrors.DomainError) {
	var version models.ContentVersion
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&version)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "バージョンが見つかりません")
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return &version, nil
}

func (r *VersionRepositoryImpl) FindByContentID(ctx context.Context, contentID string) ([]*models.ContentVersion, *myerrors.DomainError) {
	var versions []*models.ContentVersion
	result := r.db.WithContext(ctx).Where("content_id = ?", contentID).Find(&versions)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return versions, nil
}

func (r *VersionRepositoryImpl) FindLatestByContentID(ctx context.Context, contentID string) (*models.ContentVersion, *myerrors.DomainError) {
	var version models.ContentVersion
	// Note: Assumes version field is numeric and sortable for proper ordering
	result := r.db.WithContext(ctx).Where("content_id = ?", contentID).Order("version desc").First(&version)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "バージョンが見つかりません")
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return &version, nil
}

func (r *VersionRepositoryImpl) Delete(ctx context.Context, id string) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.ContentVersion{})
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	if result.RowsAffected == 0 {
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "バージョンが見つかりません")
	}
	return nil
}
