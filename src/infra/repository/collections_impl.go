package infrastructure

import (
	"errors"
	"fmt"

	"w3st/domain/models"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type CollectionsRepository struct {
	db *gorm.DB
}

func NewCollectionsRepository(db *gorm.DB) *CollectionsRepository {
	return &CollectionsRepository{
		db: db,
	}
}

func (r *CollectionsRepository) CreateCollection(newCollection *models.ApiCollection) (collection error) {
	result := r.db.Create(newCollection)

	if result.Error != nil {
		// クエリの実行中に発生したエラー
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	// コレクションの作成に成功した場合
	return nil
}

func (r *CollectionsRepository) GetCollectionByProjectId(projectId int) ([]models.ApiCollection, error) {
	var collection []models.ApiCollection
	result := r.db.Where("project_id = ?", projectId).Find(&collection)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "コレクションが見つかりません")
		}
		// その他のエラー
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	// 中身を出力
	fmt.Print(collection)

	return collection, nil
}

func (r *CollectionsRepository) GetCollectionsByCollectionId(collectionId int, projectId int) (*models.ApiCollection, error) {
	var collection models.ApiCollection
	result := r.db.Where("id = ? AND project_id = ?", collectionId, projectId).First(&collection)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "コレクションが見つかりません")
		}
		// その他のエラー
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return &collection, nil
}
