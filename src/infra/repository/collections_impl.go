package infrastructure

import (
	"gorm.io/gorm"
	"w3st/domain/models"
	myerrors "w3st/errors"
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
		return myerrors.NewDomainError(myerrors.QueryError, result.Error.Error())
	}

	// コレクションの作成に成功した場合
	return nil

}
