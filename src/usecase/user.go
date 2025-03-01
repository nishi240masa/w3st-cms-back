package usecase

import (
	"w3st/interfaces/repositories"
	"w3st/interfaces/services"
	"w3st/models"
	"w3st/presenter"
	"w3st/utils"
)

type UserUsecase interface {
	Create(newUser *models.User) (string, error)
	FindByEmail(email string) (string, error)
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

func (u *userUsecase) Create(newUser *models.User) (string, error) {

	err := u.userRepo.Create(newUser)
	if err != nil {
		return "", err
	}
	stringID := utils.UuidToString(newUser.ID)
	
	token, err := u.authService.GenerateToken(stringID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) FindByEmail(email string) (string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	response := u.userPresenter.ResponseUser(user)
	token, err := u.authService.GenerateToken(response.ID.String())
	if err != nil {
		return "", err
	}
	return token, nil
}