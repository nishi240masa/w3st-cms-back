package usecase

import (
	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"github.com/google/uuid"
)

type FieldUsecase interface {
	Create(projectId int, newField *models.FieldData) error
	Update(projectId int, newField *models.FieldData) error
	Delete(projectId int, fieldId string) error
	GetByCollectionId(collectionId int, projectId int) ([]models.FieldData, error)
}

type fieldUsecase struct {
	fieldRepo      repositories.FieldRepository
	collectionRepo repositories.CollectionsRepository
}

func NewFieldUsecase(fieldRepo repositories.FieldRepository, collectionRepo repositories.CollectionsRepository) FieldUsecase {
	return &fieldUsecase{
		fieldRepo:      fieldRepo,
		collectionRepo: collectionRepo,
	}
}

func (f *fieldUsecase) Create(projectId int, newField *models.FieldData) error {
	// collectionが存在するか確認
	if _, err := f.collectionRepo.GetCollectionsByCollectionId(newField.CollectionID, projectId); err != nil {
		//	collectionが存在しない場合
		return myerrors.WrapDomainError("fieldUsecase.Create", err)
	}
	if err := f.fieldRepo.CreateField(newField); err != nil {
		// フィールドの作成に失敗した場合
		return myerrors.WrapDomainError("fieldUsecase.Create", err)
	}
	// フィールドの作成に成功した場合
	return nil
}

func (f *fieldUsecase) GetByCollectionId(collectionId int, projectId int) ([]models.FieldData, error) {
	// collectionが存在するか確認
	if _, err := f.collectionRepo.GetCollectionsByCollectionId(collectionId, projectId); err != nil {
		return nil, myerrors.WrapDomainError("fieldUsecase.GetByCollectionId", err)
	}

	fields, err := f.fieldRepo.GetFieldsByCollectionId(collectionId, projectId)
	if err != nil {
		return nil, myerrors.WrapDomainError("fieldUsecase.GetByCollectionId", err)
	}
	return fields, nil
}

func (f *fieldUsecase) Update(projectId int, newField *models.FieldData) error {
	// collectionが存在するか確認
	if _, err := f.collectionRepo.GetCollectionsByCollectionId(newField.CollectionID, projectId); err != nil {
		return myerrors.WrapDomainError("fieldUsecase.Update", err)
	}
	if err := f.fieldRepo.UpdateField(newField); err != nil {
		// フィールドの更新に失敗した場合
		return myerrors.WrapDomainError("fieldUsecase.Update", err)
	}
	// フィールドの更新に成功した場合
	return nil
}

func (f *fieldUsecase) Delete(projectId int, fieldId string) error {
	fieldUuid, err := uuid.Parse(fieldId)
	if err != nil {
		return myerrors.WrapDomainError("fieldUsecase.Delete", err)
	}
	// フィールドを削除する
	err = f.fieldRepo.DeleteFieldById(projectId, fieldUuid)
	if err != nil {
		return myerrors.WrapDomainError("fieldUsecase.Delete", err)
	}
	return nil
}
