package presenter

import (
	"w3st/domain/models"
	"w3st/dto"
)

type PermissionPresenter interface {
	ResponsePermission(permission *models.UserPermission) *dto.PermissionResponse
	ResponsePermissions(permissions []*models.UserPermission) []*dto.PermissionResponse
}

type permissionPresenter struct{}

func NewPermissionPresenter() PermissionPresenter {
	return &permissionPresenter{}
}

func (p *permissionPresenter) ResponsePermission(permission *models.UserPermission) *dto.PermissionResponse {
	return &dto.PermissionResponse{
		ID:         permission.ID.String(),
		UserID:     permission.UserID.String(),
		Permission: permission.Permission,
		Resource:   permission.Resource,
		CreatedAt:  permission.CreatedAt.Format(ISO8601Format),
		UpdatedAt:  permission.UpdatedAt.Format(ISO8601Format),
	}
}

func (p *permissionPresenter) ResponsePermissions(permissions []*models.UserPermission) []*dto.PermissionResponse {
	responses := make([]*dto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responses[i] = p.ResponsePermission(permission)
	}
	return responses
}
