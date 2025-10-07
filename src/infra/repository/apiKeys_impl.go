package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type ApiKeyRepositoryImpl struct {
	db *gorm.DB
}

func NewApiKeyRepositoryImpl(db *gorm.DB) repositories.ApiKeyRepository {
	return &ApiKeyRepositoryImpl{db: db}
}

func (r *ApiKeyRepositoryImpl) Create(ctx context.Context, apiKey *models.ApiKeys) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Create(apiKey)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return nil
}

func (r *ApiKeyRepositoryImpl) FindByKey(ctx context.Context, key string) (*models.ApiKeys, *myerrors.DomainError) {
	var apiKey models.ApiKeys
	result := r.db.WithContext(ctx).Where("key = ? AND revoked = false", key).First(&apiKey)
	if result.Error != nil {
		// sqlmock may return sql.ErrNoRows while gorm uses gorm.ErrRecordNotFound in some setups.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, sql.ErrNoRows) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "APIキーが見つかりません")
		}
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return &apiKey, nil
}

func (r *ApiKeyRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*models.ApiKeys, *myerrors.DomainError) {
	var apiKeys []*models.ApiKeys
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&apiKeys)
	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return apiKeys, nil
}

func (r *ApiKeyRepositoryImpl) Update(ctx context.Context, apiKey *models.ApiKeys) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Save(apiKey)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return nil
}

func (r *ApiKeyRepositoryImpl) Delete(ctx context.Context, id int) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.ApiKeys{})
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	if result.RowsAffected == 0 {
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "APIキーが見つかりません")
	}
	return nil
}