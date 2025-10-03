package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type VersionUsecase interface {
	CreateVersion(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, data interface{}) (*models.ContentVersion, error)
	GetVersionsByContentID(ctx context.Context, userID uuid.UUID, contentID uuid.UUID) ([]*models.ContentVersion, error)
	GetLatestVersion(ctx context.Context, userID uuid.UUID, contentID uuid.UUID) (*models.ContentVersion, error)
	RestoreVersion(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, versionID uuid.UUID) (*models.ContentVersion, error)
}

type versionUsecase struct {
	versionRepo repositories.VersionRepository
}

func NewVersionUsecase(versionRepo repositories.VersionRepository) VersionUsecase {
	return &versionUsecase{
		versionRepo: versionRepo,
	}
}

func (v *versionUsecase) CreateVersion(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, data interface{}) (*models.ContentVersion, error) {
	// 最新バージョンを取得してバージョン番号を決定
	latest, err := v.versionRepo.FindLatestByContentID(ctx, contentID.String())
	version := 1
	if err == nil && latest != nil {
		version = latest.Version + 1
	} else if err != nil {
		// エラーが QueryDataNotFoundError 以外ならエラー
		var domainErr *myerrors.DomainError
		if !errors.As(err, &domainErr) || domainErr.GetType() != myerrors.QueryDataNotFoundError {
			return nil, myerrors.WrapDomainError("versionUsecase.CreateVersion", err)
		}
		// Not found は version = 1 でOK
	}

	// JSON データに変換
	jsonBytes, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.InvalidParameter, "データのJSON変換に失敗しました")
	}
	jsonData := datatypes.JSON(jsonBytes)

	// ContentVersion を作成
	contentVersion := &models.ContentVersion{
		ContentID: contentID,
		Version:   version,
		Data:      jsonData,
		UserID:    userID,
	}

	// リポジトリで作成
	if err := v.versionRepo.Create(ctx, contentVersion); err != nil {
		return nil, myerrors.WrapDomainError("versionUsecase.CreateVersion", err)
	}

	return contentVersion, nil
}

func (v *versionUsecase) GetVersionsByContentID(ctx context.Context, userID uuid.UUID, contentID uuid.UUID) ([]*models.ContentVersion, error) {
	versions, err := v.versionRepo.FindByContentID(ctx, contentID.String())
	if err != nil {
		return nil, myerrors.WrapDomainError("versionUsecase.GetVersionsByContentID", err)
	}

	// 所有者チェック (最初のバージョンのUserIDでチェック)
	if len(versions) > 0 && versions[0].UserID.String() != userID.String() {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.UnPermittedOperation, "アクセス権限がありません")
	}

	return versions, nil
}

func (v *versionUsecase) GetLatestVersion(ctx context.Context, userID uuid.UUID, contentID uuid.UUID) (*models.ContentVersion, error) {
	latest, err := v.versionRepo.FindLatestByContentID(ctx, contentID.String())
	if err != nil {
		return nil, myerrors.WrapDomainError("versionUsecase.GetLatestVersion", err)
	}

	// 所有者チェック
	if latest.UserID.String() != userID.String() {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.UnPermittedOperation, "アクセス権限がありません")
	}

	return latest, nil
}

func (v *versionUsecase) RestoreVersion(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, versionID uuid.UUID) (*models.ContentVersion, error) {
	// 指定バージョンを取得
	version, err := v.versionRepo.FindByID(ctx, versionID.String())
	if err != nil {
		return nil, myerrors.WrapDomainError("versionUsecase.RestoreVersion", err)
	}

	// 所有者チェック
	if version.UserID.String() != userID.String() {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.UnPermittedOperation, "アクセス権限がありません")
	}

	// ContentID の一致チェック
	if version.ContentID.String() != contentID.String() {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.InvalidParameter, "バージョンが指定されたコンテンツに属していません")
	}

	return version, nil
}
