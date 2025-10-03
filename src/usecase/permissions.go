package usecase

import (
	"context"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"github.com/google/uuid"
)

type PermissionUsecase interface {
	CheckPermission(ctx context.Context, userID uuid.UUID, permission, resource string) (bool, error)
	GrantPermission(ctx context.Context, userID uuid.UUID, permission, resource string) error
	RevokePermission(ctx context.Context, userID uuid.UUID, permission, resource string) error
	GetPermissionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPermission, error)
}

type permissionUsecase struct {
	permissionRepo repositories.PermissionRepository
}

func NewPermissionUsecase(permissionRepo repositories.PermissionRepository) PermissionUsecase {
	return &permissionUsecase{
		permissionRepo: permissionRepo,
	}
}

func (p *permissionUsecase) CheckPermission(ctx context.Context, userID uuid.UUID, permission, resource string) (bool, error) {
	permissions, err := p.permissionRepo.FindByUserIDAndResource(ctx, userID.String(), resource)
	if err != nil {
		return false, myerrors.WrapDomainError("permissionUsecase.CheckPermission", err)
	}

	for _, perm := range permissions {
		if perm.Permission == permission {
			return true, nil
		}
	}

	return false, nil
}

func (p *permissionUsecase) GrantPermission(ctx context.Context, userID uuid.UUID, permission, resource string) error {
	// すでに権限があるかチェック
	hasPermission, err := p.CheckPermission(ctx, userID, permission, resource)
	if err != nil {
		return myerrors.WrapDomainError("permissionUsecase.GrantPermission", err)
	}
	if hasPermission {
		return myerrors.NewDomainErrorWithMessage(myerrors.AlreadyExist, "すでに権限が付与されています")
	}

	// UserPermission を作成
	userPerm := &models.UserPermission{
		UserID:     userID,
		Permission: permission,
		Resource:   resource,
	}

	// リポジトリで作成
	if err := p.permissionRepo.Create(ctx, userPerm); err != nil {
		return myerrors.WrapDomainError("permissionUsecase.GrantPermission", err)
	}

	return nil
}

func (p *permissionUsecase) RevokePermission(ctx context.Context, userID uuid.UUID, permission, resource string) error {
	// 権限を取得
	permissions, err := p.permissionRepo.FindByUserIDAndResource(ctx, userID.String(), resource)
	if err != nil {
		return myerrors.WrapDomainError("permissionUsecase.RevokePermission", err)
	}

	var permID string
	for _, perm := range permissions {
		if perm.Permission == permission {
			permID = perm.ID.String()
			break
		}
	}

	if permID == "" {
		return myerrors.NewDomainErrorWithMessage(myerrors.QueryDataNotFoundError, "権限が見つかりません")
	}

	// 削除
	if err := p.permissionRepo.Delete(ctx, permID); err != nil {
		return myerrors.WrapDomainError("permissionUsecase.RevokePermission", err)
	}

	return nil
}

func (p *permissionUsecase) GetPermissionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.UserPermission, error) {
	permissions, err := p.permissionRepo.FindByUserID(ctx, userID.String())
	if err != nil {
		return nil, myerrors.WrapDomainError("permissionUsecase.GetPermissionsByUser", err)
	}

	return permissions, nil
}
