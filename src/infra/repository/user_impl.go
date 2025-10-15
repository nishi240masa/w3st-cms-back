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

		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
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
			return &user, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "ユーザーが見つかりません")
		}
		// その他のエラー
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return &user, nil
}

// GetAllUsers
func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]models.Users, *myerrors.DomainError) {
	var users []models.Users
	result := r.db.WithContext(ctx).Find(&users)

	if result.Error != nil {
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return users, nil
}

// DeleteUser
func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, userID string) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Where("id = ?", userID).Delete(&models.Users{})

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	if result.RowsAffected == 0 {
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "ユーザーが見つかりません")
	}

	return nil
}

// Update
func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.Users) *myerrors.DomainError {
	result := r.db.WithContext(ctx).Save(user)

	if result.Error != nil {
		return myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}
	return nil
}

// FindByID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, userID string) (*models.Users, *myerrors.DomainError) {
	var user models.Users
	result := r.db.WithContext(ctx).Where("id = ?", userID).First(&user)

	// エラーが発生した場合
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &user, myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "ユーザーが見つかりません")
		}
		// その他のエラー
		fmt.Println("その他のエラー:", result.Error)
		return nil, myerrors.NewDomainError(myerrors.QueryError, result.Error)
	}

	return &user, nil
}
