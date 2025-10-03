package infrastructure

import (
	"context"
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type MediaRepositoryImpl struct {
	db *gorm.DB
}

func NewMediaRepositoryImpl(db *gorm.DB) repositories.MediaRepository {
	return &MediaRepositoryImpl{db: db}
}

func (r *MediaRepositoryImpl) Create(ctx context.Context, media *models.MediaAsset) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Create(media)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return nil
}

func (r *MediaRepositoryImpl) FindByID(ctx context.Context, id string) (*models.MediaAsset, *myerrors.DomainError) {
	var media models.MediaAsset
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&media)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "メディアが見つかりません")
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return &media, nil
}

func (r *MediaRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*models.MediaAsset, *myerrors.DomainError) {
	var medias []*models.MediaAsset
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&medias)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return medias, nil
}

func (r *MediaRepositoryImpl) Update(ctx context.Context, media *models.MediaAsset) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Model(&models.MediaAsset{}).Where("id = ?", media.ID).Updates(media)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	if result.RowsAffected == 0 {
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "メディアが見つかりません")
	}
	return nil
}

func (r *MediaRepositoryImpl) Delete(ctx context.Context, id string) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MediaAsset{})
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	if result.RowsAffected == 0 {
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "メディアが見つかりません")
	}
	return nil
}
