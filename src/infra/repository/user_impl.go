package infrastructure

import (
	"context"
	"w3st/domain/models"
	"w3st/domain/repositories"
	"w3st/errors"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
    db *gorm.DB
}


func NewUserRepositoryImpl(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}



// Create
func (r *UserRepositoryImpl) Create(ctx context.Context, newUser *models.User)  *errors.DomainError { // context を引数に追加
	result := r.db.WithContext(ctx).Create(newUser)
	if result.Error != nil {
		// クエリの実行中に発生したエラー
		return errors.NewDomainError(errors.QueryError, result.Error.Error())
	}
	// ユーザーの作成に成功した場合
	return nil
}

// Find
func (r *UserRepositoryImpl) FindByEmail (ctx context.Context,email string) (*models.User,  *errors.DomainError) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	// エラーが発生した場合
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewDomainError(errors.QueryDataNotFoundError, "ユーザーが見つかりません")
		}
		// その他のエラー
		return nil, errors.NewDomainError(errors.QueryError, result.Error.Error())
	}

	// ユーザーが見つからなかった場合
	if result.RowsAffected == 0 {
		return nil, errors.NewDomainError(errors.QueryDataNotFoundError, "ユーザーが見つかりません")
	}
	
	return &user, nil
	
}

