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
	InitApiKeyUsecase() usecase.ApiKeyUsecase
	InitApiKeyController() *controllers.ApiKeyController
	InitSDKCollectionsController() *controllers.SDKCollectionsController
	InitGUICollectionsController() *controllers.GUICollectionsController
	InitSDKEntriesController() *controllers.SDKEntriesController
	InitGUIEntriesController() *controllers.GUIEntriesController
	InitFieldController() *controllers.FieldController
	InitMediaController() *controllers.MediaController
	InitAuditController() *controllers.AuditController
	InitPermissionController() *controllers.PermissionController
	InitVersionController() *controllers.VersionController
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

func (f factory) InitApiKeyUsecase() usecase.ApiKeyUsecase {
	return usecase.NewApiKeyUsecase()
}

func (f factory) InitApiKeyController() *controllers.ApiKeyController {
	apiKeyUsecase := usecase.NewApiKeyUsecase()

	return controllers.NewApiKeyController(apiKeyUsecase)
}

func (f factory) InitSDKCollectionsController() *controllers.SDKCollectionsController {
	collectionRepo := infrastructure.NewCollectionsRepository(f.DB)
	collectionUsecase := usecase.NewCollectionsUsecase(collectionRepo)

	return controllers.NewSDKCollectionsController(collectionUsecase)
}

func (f factory) InitGUICollectionsController() *controllers.GUICollectionsController {
	collectionRepo := infrastructure.NewCollectionsRepository(f.DB)
	fieldRepo := infrastructure.NewFieldRepository(f.DB)
	collectionUsecase := usecase.NewCollectionsUsecase(collectionRepo)
	fieldUsecase := usecase.NewFieldUsecase(fieldRepo, collectionRepo)

	return controllers.NewGUICollectionsController(collectionUsecase, fieldUsecase)
}

func (f factory) InitSDKEntriesController() *controllers.SDKEntriesController {
	entriesRepo := infrastructure.NewEntriesRepository(f.DB)
	collectionRepo := infrastructure.NewCollectionsRepository(f.DB)
	collectionUsecase := usecase.NewCollectionsUsecase(collectionRepo)
	entriesUsecase := usecase.NewEntriesUsecase(entriesRepo, collectionUsecase)

	return controllers.NewSDKEntriesController(entriesUsecase)
}

func (f factory) InitGUIEntriesController() *controllers.GUIEntriesController {
	entriesRepo := infrastructure.NewEntriesRepository(f.DB)
	collectionRepo := infrastructure.NewCollectionsRepository(f.DB)
	collectionUsecase := usecase.NewCollectionsUsecase(collectionRepo)
	entriesUsecase := usecase.NewEntriesUsecase(entriesRepo, collectionUsecase)

	return controllers.NewGUIEntriesController(entriesUsecase)
}

func (f factory) InitFieldController() *controllers.FieldController {
	fieldRepo := infrastructure.NewFieldRepository(f.DB)
	collectionRepo := infrastructure.NewCollectionsRepository(f.DB)
	fieldUsecase := usecase.NewFieldUsecase(fieldRepo, collectionRepo)

	return controllers.NewFieldController(fieldUsecase)
}

func (f factory) InitMediaController() *controllers.MediaController {
	mediaRepo := infrastructure.NewMediaRepositoryImpl(f.DB)
	mediaUsecase := usecase.NewMediaUsecase(mediaRepo)

	return controllers.NewMediaController(mediaUsecase)
}

func (f factory) InitAuditController() *controllers.AuditController {
	auditRepo := infrastructure.NewAuditRepositoryImpl(f.DB)
	auditUsecase := usecase.NewAuditUsecase(auditRepo)

	return controllers.NewAuditController(auditUsecase)
}

func (f factory) InitPermissionController() *controllers.PermissionController {
	permissionRepo := infrastructure.NewPermissionRepositoryImpl(f.DB)
	permissionUsecase := usecase.NewPermissionUsecase(permissionRepo)

	return controllers.NewPermissionController(permissionUsecase)
}

func (f factory) InitVersionController() *controllers.VersionController {
	versionRepo := infrastructure.NewVersionRepositoryImpl(f.DB)
	versionUsecase := usecase.NewVersionUsecase(versionRepo)

	return controllers.NewVersionController(versionUsecase)
}
