package usecase

import (
	"context" // context をインポート
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
)

type UserUsecase interface {
	Create(newUser *models.Users, ctx context.Context) (*models.Users, error)
	FindByEmail(email string) (*models.Users, error)
	FindByID(userID string) (*models.Users, error)
}

type userUsecase struct {
	userRepo repositories.UserRepository
}

func NewUserUsecase(userRepo repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) Create(newUser *models.Users, ctx context.Context) (*models.Users, error) {
	// すでに存在するか確認
	_, err := u.userRepo.FindByEmail(ctx, newUser.Email)
	if err != nil {
		// ユーザーが存在しない場合
		if errors.Is(err, &myerrors.DomainError{ErrType: myerrors.QueryDataNotFoundError}) {
			if err := u.userRepo.Create(ctx, newUser); err != nil {
				return nil, myerrors.WrapDomainError("usecase.Create", err)
			}
			return newUser, nil
		}
		// それ以外のエラー（DB障害など）
		return nil, myerrors.WrapDomainError("usecase.Create", err)
	}

	// ユーザーがすでに存在していた場合
	return nil, myerrors.NewDomainErrorWithMessage(myerrors.AlreadyExist, "ユーザーはすでに存在します")
}

func (u *userUsecase) FindByEmail(email string) (*models.Users, error) {
	user, err := u.userRepo.FindByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) FindByID(userID string) (*models.Users, error) {
	user, err := u.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
