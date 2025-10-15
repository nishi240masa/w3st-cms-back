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
	Update(user *models.Users, ctx context.Context) error
	GetAllUsers() ([]models.Users, error)
	DeleteUser(userID string) error
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

func (u *userUsecase) Update(user *models.Users, ctx context.Context) error {
	if err := u.userRepo.Update(ctx, user); err != nil {
		return myerrors.WrapDomainError("usecase.Update", err)
	}
	return nil
}

func (u *userUsecase) GetAllUsers() ([]models.Users, error) {
	users, err := u.userRepo.GetAllUsers(context.Background())
	if err != nil {
		return nil, myerrors.WrapDomainError("usecase.GetAllUsers", err)
	}
	return users, nil
}

func (u *userUsecase) DeleteUser(userID string) error {
	if err := u.userRepo.DeleteUser(context.Background(), userID); err != nil {
		return myerrors.WrapDomainError("usecase.DeleteUser", err)
	}
	return nil
}
