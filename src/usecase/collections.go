package usecase

import (
	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
)

type CollectionsUsecase interface {
	Make(newCollection *models.ApiCollection) error
}

type collectionsUsecase struct {
	collectionsRepo repositories.CollectionsRepository
}

func NewCollectionsUsecase(collectionsRepo repositories.CollectionsRepository) CollectionsUsecase {
	return &collectionsUsecase{
		collectionsRepo: collectionsRepo,
	}
}

func (c *collectionsUsecase) Make(newCollection *models.ApiCollection) error {

	// コレクションを作成する
	err := c.collectionsRepo.CreateCollection(newCollection)
	if err != nil {
		// エラー処理
		return myerrors.NewDomainError(myerrors.QueryError, err.Error())
	}
	return nil

}
