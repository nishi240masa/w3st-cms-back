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
	InitCollectionController() *controllers.CollectionsController
	InitFieldController() *controllers.FieldController
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

func (f factory) InitCollectionController() *controllers.CollectionsController {
	collectionRepo := infrastructure.NewCollectionsRepository(f.DB)
	collectionUsecase := usecase.NewCollectionsUsecase(collectionRepo)

	return controllers.NewCollectionsController(collectionUsecase)
}

func (f factory) InitFieldController() *controllers.FieldController {
	fieldRepo := infrastructure.NewFieldRepository(f.DB)
	collectionRepo := infrastructure.NewCollectionsRepository(f.DB)
	fieldUsecase := usecase.NewFieldUsecase(fieldRepo, collectionRepo)

	return controllers.NewFieldController(fieldUsecase)
}
