package usecase

import (
	"context"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"

	"github.com/google/uuid"
)

type AuditUsecase interface {
	LogAction(ctx context.Context, userID uuid.UUID, action, resource, details string) error
	LogActionWithProject(ctx context.Context, userID uuid.UUID, projectID int, action, resource, details string) error
	GetLogsByUser(ctx context.Context, userID uuid.UUID) ([]*models.AuditLog, error)
	GetLogsByProject(ctx context.Context, projectID int) ([]*models.AuditLog, error)
	GetLogsByAction(ctx context.Context, action string) ([]*models.AuditLog, error)
	GetAllLogs(ctx context.Context, limit int, offset int) ([]*models.AuditLog, error)
	GetLogsByProjectWithLimit(ctx context.Context, projectID int, limit int, offset int) ([]*models.AuditLog, error)
}

type auditUsecase struct {
	auditRepo repositories.AuditRepository
}

func NewAuditUsecase(auditRepo repositories.AuditRepository) AuditUsecase {
	return &auditUsecase{
		auditRepo: auditRepo,
	}
}

func (a *auditUsecase) LogAction(ctx context.Context, userID uuid.UUID, action, resource, details string) error {
	// デフォルトプロジェクトIDを使用（後方互換性のため）
	return a.LogActionWithProject(ctx, userID, 1, action, resource, details)
}

func (a *auditUsecase) LogActionWithProject(ctx context.Context, userID uuid.UUID, projectID int, action, resource, details string) error {
	// AuditLog を作成
	log := &models.AuditLog{
		UserID:    userID,
		ProjectID: projectID,
		Action:    action,
		Resource:  resource,
		Details:   details,
	}

	// リポジトリで作成
	if err := a.auditRepo.Create(ctx, log); err != nil {
		return myerrors.WrapDomainError("auditUsecase.LogActionWithProject", err)
	}

	return nil
}

func (a *auditUsecase) GetLogsByUser(ctx context.Context, userID uuid.UUID) ([]*models.AuditLog, error) {
	logs, err := a.auditRepo.FindByUserID(ctx, userID.String())
	if err != nil {
		return nil, myerrors.WrapDomainError("auditUsecase.GetLogsByUser", err)
	}

	return logs, nil
}

func (a *auditUsecase) GetLogsByAction(ctx context.Context, action string) ([]*models.AuditLog, error) {
	logs, err := a.auditRepo.FindByAction(ctx, action)
	if err != nil {
		return nil, myerrors.WrapDomainError("auditUsecase.GetLogsByAction", err)
	}

	return logs, nil
}

func (a *auditUsecase) GetLogsByProject(ctx context.Context, projectID int) ([]*models.AuditLog, error) {
	logs, err := a.auditRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, myerrors.WrapDomainError("auditUsecase.GetLogsByProject", err)
	}

	return logs, nil
}

func (a *auditUsecase) GetLogsByProjectWithLimit(ctx context.Context, projectID int, limit int, offset int) ([]*models.AuditLog, error) {
	logs, err := a.auditRepo.FindByProjectIDWithLimit(ctx, projectID, limit, offset)
	if err != nil {
		return nil, myerrors.WrapDomainError("auditUsecase.GetLogsByProjectWithLimit", err)
	}

	return logs, nil
}

func (a *auditUsecase) GetAllLogs(ctx context.Context, limit int, offset int) ([]*models.AuditLog, error) {
	logs, err := a.auditRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, myerrors.WrapDomainError("auditUsecase.GetAllLogs", err)
	}

	return logs, nil
}
