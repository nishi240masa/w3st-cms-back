package presenter

import (
	"w3st/domain/models"
	"w3st/dto"
)

type VersionPresenter interface {
	ResponseVersion(version *models.ContentVersion) *dto.VersionResponse
	ResponseVersions(versions []*models.ContentVersion) []*dto.VersionResponse
}

type versionPresenter struct{}

func NewVersionPresenter() VersionPresenter {
	return &versionPresenter{}
}

func (v *versionPresenter) ResponseVersion(version *models.ContentVersion) *dto.VersionResponse {
	return &dto.VersionResponse{
		ID:        version.ID.String(),
		ContentID: version.ContentID.String(),
		Version:   version.Version,
		Data:      string(version.Data),
		UserID:    version.UserID.String(),
		CreatedAt: version.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: version.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (v *versionPresenter) ResponseVersions(versions []*models.ContentVersion) []*dto.VersionResponse {
	responses := make([]*dto.VersionResponse, len(versions))
	for i, version := range versions {
		responses[i] = v.ResponseVersion(version)
	}
	return responses
}
