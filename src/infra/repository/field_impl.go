package infrastructure

import (
	"w3st/domain/models"
	myerrors "w3st/errors"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type FieldRepository struct {
	db *gorm.DB
}

func NewFieldRepository(db *gorm.DB) *FieldRepository {
	return &FieldRepository{
		db: db,
	}
}

func (r *FieldRepository) CreateField(newField *models.FieldData) error {
	result := r.db.Create(newField)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	// フィールドの作成に成功した場合
	return nil
}

func (r *FieldRepository) UpdateField(newField *models.FieldData) error {
	result := r.db.Save(newField)
	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	// フィールドの更新に成功した場合
	return nil
}

func (r *FieldRepository) DeleteFieldById(userId uuid.UUID, fieldId uuid.UUID) error {
	result := r.db.Where("user_id = ? AND id = ?", userId.String(), fieldId.String()).Delete(&models.FieldData{})
	if result.Error != nil {
		// クエリの実行中に発生したエラー
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	if result.RowsAffected == 0 {
		// レコードが見つからなかった場合
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "フィールドが見つかりません")
	}
	return nil
}
