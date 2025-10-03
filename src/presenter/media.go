package presenter

import (
	"w3st/domain/models"
	"w3st/dto"
)

type MediaPresenter interface {
	ResponseMedia(media *models.MediaAsset) *dto.MediaResponse
	ResponseMedias(medias []*models.MediaAsset) []*dto.MediaResponse
}

type mediaPresenter struct{}

func NewMediaPresenter() MediaPresenter {
	return &mediaPresenter{}
}

func (m *mediaPresenter) ResponseMedia(media *models.MediaAsset) *dto.MediaResponse {
	return &dto.MediaResponse{
		ID:        media.ID.String(),
		Name:      media.Name,
		Type:      media.Type,
		Path:      media.Path,
		Size:      media.Size,
		UserID:    media.UserID.String(),
		CreatedAt: media.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: media.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (m *mediaPresenter) ResponseMedias(medias []*models.MediaAsset) []*dto.MediaResponse {
	responses := make([]*dto.MediaResponse, len(medias))
	for i, media := range medias {
		responses[i] = m.ResponseMedia(media)
	}
	return responses
}
