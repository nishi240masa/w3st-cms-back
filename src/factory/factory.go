package factory

import (
	infrastructure "w3st/infra/repository"
	"w3st/interfaces/controllers"
	"w3st/presenter"
	"w3st/usecase"

	"gorm.io/gorm"
)

type Factory interface {
	InitUserController() *controllers.UserController
	InitAuthUsecase() usecase.JwtUsecase
}

type factory struct {
	DB *gorm.DB
}

func NewFactory(db *gorm.DB) Factory {
	return &factory{DB: db}
}

func (f factory) InitUserController() *controllers.UserController {
	userRepo := infrastructure.NewUserRepositoryImpl(f.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	jwtAuthUsecase := usecase.NewjwtAuthUsecase()
	userPresenter := presenter.NewUserPresenter()

	return controllers.NewUserController(userUsecase, jwtAuthUsecase, userPresenter)
}

func (f factory) InitAuthUsecase() usecase.JwtUsecase {
	return usecase.NewjwtAuthUsecase()
}
