package factory

import (
	infrastructure "w3st/infra/repository"
	"w3st/interfaces/controllers"
	"w3st/interfaces/services"
	"w3st/presenter"
	"w3st/usecase"

	"gorm.io/gorm"
)

type Factory interface {
	InitUserController() *controllers.UserController
	InitAuthService() services.AuthService
}

type factory struct {
	DB *gorm.DB
}

func NewFactory(db *gorm.DB) Factory {
	return &factory{DB: db}
}

func (f factory) InitUserController() *controllers.UserController {
	userRepo := infrastructure.NewUserRepositoryImpl(f.DB)
	authService := services.NewAuthService()
	txRepo := infrastructure.NewTransactionRepositoryImpl(f.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, authService, txRepo)
	userPresenter := presenter.NewUserPresenter()
	return controllers.NewUserController(userUsecase, authService, userPresenter)
}

func (f factory) InitAuthService() services.AuthService {
	return services.NewAuthService()
}
