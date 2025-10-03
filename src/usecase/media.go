package usecase

import (
	"context"
	"path/filepath"
	"strings"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"github.com/google/uuid"
)

type MediaUsecase interface {
	Upload(ctx context.Context, userID uuid.UUID, name, fileType, path string, size int64) (*models.MediaAsset, error)
	GetByID(ctx context.Context, userID uuid.UUID, id string) (*models.MediaAsset, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.MediaAsset, error)
	Delete(ctx context.Context, userID uuid.UUID, id string) error
}

type mediaUsecase struct {
	mediaRepo repositories.MediaRepository
}

func NewMediaUsecase(mediaRepo repositories.MediaRepository) MediaUsecase {
	return &mediaUsecase{
		mediaRepo: mediaRepo,
	}
}

func (m *mediaUsecase) Upload(ctx context.Context, userID uuid.UUID, name, fileType, path string, size int64) (*models.MediaAsset, error) {
	// ファイルサイズのチェック (例: 10MB 制限)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if size > maxSize {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.InvalidParameter, "ファイルサイズが大きすぎます")
	}

	// ファイルタイプのチェック
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif", "application/pdf"}
	isAllowed := false
	for _, allowed := range allowedTypes {
		if fileType == allowed {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.InvalidParameter, "許可されていないファイルタイプです")
	}

	// ファイル拡張子のチェック
	ext := strings.ToLower(filepath.Ext(name))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"}
	isAllowedExt := false
	for _, allowed := range allowedExts {
		if ext == allowed {
			isAllowedExt = true
			break
		}
	}
	if !isAllowedExt {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.InvalidParameter, "許可されていないファイル拡張子です")
	}

	// MediaAsset を作成
	media := &models.MediaAsset{
		Name:   name,
		Type:   fileType,
		Path:   path,
		Size:   size,
		UserID: userID,
	}

	// リポジトリで作成
	if err := m.mediaRepo.Create(ctx, media); err != nil {
		return nil, myerrors.WrapDomainError("mediaUsecase.Upload", err)
	}

	return media, nil
}

func (m *mediaUsecase) GetByID(ctx context.Context, userID uuid.UUID, id string) (*models.MediaAsset, error) {
	media, err := m.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, myerrors.WrapDomainError("mediaUsecase.GetByID", err)
	}

	// 所有者チェック
	if media.UserID.String() != userID.String() {
		return nil, myerrors.NewDomainErrorWithMessage(myerrors.UnPermittedOperation, "アクセス権限がありません")
	}

	return media, nil
}

func (m *mediaUsecase) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.MediaAsset, error) {
	medias, err := m.mediaRepo.FindByUserID(ctx, userID.String())
	if err != nil {
		return nil, myerrors.WrapDomainError("mediaUsecase.GetByUserID", err)
	}

	return medias, nil
}

func (m *mediaUsecase) Delete(ctx context.Context, userID uuid.UUID, id string) error {
	// まず所有者チェックのために取得
	media, err := m.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return myerrors.WrapDomainError("mediaUsecase.Delete", err)
	}

	if media.UserID.String() != userID.String() {
		return myerrors.NewDomainErrorWithMessage(myerrors.UnPermittedOperation, "アクセス権限がありません")
	}

	// 削除
	if err := m.mediaRepo.Delete(ctx, id); err != nil {
		return myerrors.WrapDomainError("mediaUsecase.Delete", err)
	}

	return nil
}
