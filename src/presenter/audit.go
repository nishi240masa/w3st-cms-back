package presenter

import (
	"w3st/domain/models"
	"w3st/dto"
)

type AuditPresenter interface {
	ResponseAuditLog(log *models.AuditLog) *dto.AuditLogResponse
	ResponseAuditLogs(logs []*models.AuditLog) []*dto.AuditLogResponse
}

type auditPresenter struct{}

func NewAuditPresenter() AuditPresenter {
	return &auditPresenter{}
}

func (a *auditPresenter) ResponseAuditLog(log *models.AuditLog) *dto.AuditLogResponse {
	return &dto.AuditLogResponse{
		ID:        log.ID.String(),
		UserID:    log.UserID.String(),
		Action:    log.Action,
		Resource:  log.Resource,
		Details:   log.Details,
		CreatedAt: log.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (a *auditPresenter) ResponseAuditLogs(logs []*models.AuditLog) []*dto.AuditLogResponse {
	responses := make([]*dto.AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = a.ResponseAuditLog(log)
	}
	return responses
}
