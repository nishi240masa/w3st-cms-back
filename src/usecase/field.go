package usecase

import (
	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"github.com/google/uuid"
)

type FieldUsecase interface {
	Create(userId uuid.UUID, newField *models.FieldData) error
	Update(userId uuid.UUID, newField *models.FieldData) error
	Delete(userId string, fieldId string) error
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

func (f *fieldUsecase) Create(userUuid uuid.UUID, newField *models.FieldData) error {
	// collectionが存在するか確認
	if _, err := f.collectionRepo.GetCollectionsByCollectionId(newField.CollectionID, userUuid); err != nil {
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

func (f *fieldUsecase) Update(userId uuid.UUID, newField *models.FieldData) error {
	// collectionが存在するか確認
	if _, err := f.collectionRepo.GetCollectionsByCollectionId(newField.CollectionID, userId); err != nil {
		return myerrors.WrapDomainError("fieldUsecase.Update", err)
	}
	if err := f.fieldRepo.UpdateField(newField); err != nil {
		// フィールドの更新に失敗した場合
		return myerrors.WrapDomainError("fieldUsecase.Update", err)
	}
	// フィールドの更新に成功した場合
	return nil
}

func (f *fieldUsecase) Delete(userId string, fieldId string) error {
	// UUIDに変換
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return myerrors.WrapDomainError("fieldUsecase.Delete", err)
	}
	fieldUuid, err := uuid.Parse(fieldId)
	if err != nil {
		return myerrors.WrapDomainError("fieldUsecase.Delete", err)
	}
	// フィールドを削除する
	err = f.fieldRepo.DeleteFieldById(userUuid, fieldUuid)
	if err != nil {
		return err
	}
	return nil
}
