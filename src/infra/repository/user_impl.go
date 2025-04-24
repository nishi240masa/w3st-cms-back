package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create
func (r *UserRepositoryImpl) Create(ctx context.Context, newUser *models.Users) *myerrors.DomainError { // context を引数に追加
	result := r.db.WithContext(ctx).Create(newUser)

	if result.Error != nil {
		// クエリの実行中に発生したエラー

		return myerrors.NewDomainError(myerrors.QueryError, result.Error.Error())
	}
	// ユーザーの作成に成功した場合
	return nil
}

// Find
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.Users, *myerrors.DomainError) {
	var user models.Users
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)

	// エラーが発生した場合
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &user, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, "ユーザーが見つかりません")
		}
		// その他のエラー
		fmt.Println("その他のエラー:", result.Error)
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error.Error())
	}

	return &user, nil
}

// FindByID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, userID string) (*models.Users, *myerrors.DomainError) {
	var user models.Users
	result := r.db.WithContext(ctx).Where("id = ?", userID).First(&user)

	// エラーが発生した場合
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &user, myerrors.NewDomainError(myerrors.QueryDataNotFoundError, "ユーザーが見つかりません")
		}
		// その他のエラー
		fmt.Println("その他のエラー:", result.Error)
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error.Error())
	}

	return &user, nil
}
