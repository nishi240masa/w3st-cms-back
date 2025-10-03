package presenter

import (
	"unicode/utf8"

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
	dataStr := string(version.Data)
	if !utf8.ValidString(dataStr) {
		dataStr = "" // Invalid UTF-8, set to empty string
	}
	return &dto.VersionResponse{
		ID:        version.ID.String(),
		ContentID: version.ContentID.String(),
		Version:   version.Version,
		Data:      dataStr,
		UserID:    version.UserID.String(),
		CreatedAt: version.CreatedAt.Format(ISO8601Format),
		UpdatedAt: version.UpdatedAt.Format(ISO8601Format),
	}
}

func (v *versionPresenter) ResponseVersions(versions []*models.ContentVersion) []*dto.VersionResponse {
	responses := make([]*dto.VersionResponse, len(versions))
	for i, version := range versions {
		responses[i] = v.ResponseVersion(version)
	}
	return responses
}
