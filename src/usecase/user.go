package usecase

import (
	"w3st/errors"
	"w3st/interfaces/repositories"
	"w3st/interfaces/services"
	"w3st/models"
	"w3st/presenter"
	"w3st/utils"
)

type UserUsecase interface {
	Create(newUser *models.User) (string, *errors.DomainError)
	FindByEmail(email string) (string, *errors.DomainError)
}

type userUsecase struct {
	userRepo repositories.UserRepository
	authService  services.AuthService
	userPresenter presenter.UserPresenter
}

func NewUserUsecase(userRepo repositories.UserRepository, authService services.AuthService, userPresenter presenter.UserPresenter) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		authService: authService,
		userPresenter: userPresenter,
	}
}

func (u *userUsecase) Create(newUser *models.User) (string, *errors.DomainError) {

	_, err := u.userRepo.FindByEmail(newUser.Email)
	if err == nil {
		// 既に登録されているメールアドレスの場合
		return "", errors.NewDomainError(errors.AlreadyExist, "このメールアドレスは既に登録されています")
	}


	err = u.userRepo.Create(newUser)
	if err != nil {
		// リポジトリで技術的なエラーが発生した場合
		return "", errors.NewDomainError(errors.QueryError, err.Error())

	}
	stringID := utils.UuidToString(newUser.ID)

	token, err := u.authService.GenerateToken(stringID)
	if err != nil {
		// トークン生成時にエラーが発生した場合
		return "", errors.NewDomainError(errors.ErrorUnknown, err.Error())
	}
	return token, nil
}

func (u *userUsecase) FindByEmail(email string) (string, *errors.DomainError) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		// リポジトリで技術的なエラーが発生した場合
		return "", errors.NewDomainError(errors.QueryError, err.Error())
	}
	if user == nil {
		// ユーザーが見つからなかった場合
		return "", errors.NewDomainError(errors.QueryDataNotFoundError, "ユーザーが見つかりませんでした")
	}
	
	response := u.userPresenter.ResponseUser(user)
	token, err := u.authService.GenerateToken(response.ID.String())
	if err != nil {
		// トークン生成時にエラーが発生した場合
		return "", errors.NewDomainError(errors.ErrorUnknown, err.Error())
	}
	return token, nil
}