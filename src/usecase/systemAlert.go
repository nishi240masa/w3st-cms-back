package usecase

import (
	"context"
	"time"

	"w3st/domain/models"
	"w3st/domain/repositories"
	myerrors "w3st/errors"
)

const (
	SeverityCritical = "critical"
	SeverityError    = "error"
	SeverityWarning  = "warning"
)

type SystemAlertUsecase interface {
	CreateAlert(ctx context.Context, alertType, severity, title, message string, projectID int, metadata map[string]interface{}) error
	GetActiveAlerts(ctx context.Context, projectID int) ([]models.SystemAlert, error)
	GetAllAlerts(ctx context.Context, projectID int, limit int, offset int) ([]models.SystemAlert, error)
	MarkAlertAsRead(ctx context.Context, alertID int) error
	DeleteAlert(ctx context.Context, alertID int) error
	GetAlertCount(ctx context.Context, projectID int) (int, error)
	CheckAndCreateStorageAlert(ctx context.Context, projectID int, usagePercent float64) error
	CheckAndCreateApiLimitAlert(ctx context.Context, projectID int, requestCount int, limit int) error
}

type systemAlertUsecase struct {
	systemAlertRepo repositories.SystemAlertRepository
}

func NewSystemAlertUsecase(systemAlertRepo repositories.SystemAlertRepository) SystemAlertUsecase {
	return &systemAlertUsecase{
		systemAlertRepo: systemAlertRepo,
	}
}

func (u *systemAlertUsecase) CreateAlert(ctx context.Context, alertType, severity, title, message string, projectID int, metadata map[string]interface{}) error {
	alert := &models.SystemAlert{
		AlertType: alertType,
		Severity:  severity,
		Title:     title,
		Message:   message,
		ProjectID: projectID,
		IsActive:  true,
		IsRead:    false,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := u.systemAlertRepo.Create(ctx, alert)
	if err != nil {
		return myerrors.WrapDomainError("systemAlertUsecase.CreateAlert", err)
	}

	return nil
}

func (u *systemAlertUsecase) GetActiveAlerts(ctx context.Context, projectID int) ([]models.SystemAlert, error) {
	alerts, err := u.systemAlertRepo.FindActiveByProjectID(ctx, projectID)
	if err != nil {
		return nil, myerrors.WrapDomainError("systemAlertUsecase.GetActiveAlerts", err)
	}

	return alerts, nil
}

func (u *systemAlertUsecase) GetAllAlerts(ctx context.Context, projectID int, limit int, offset int) ([]models.SystemAlert, error) {
	alerts, err := u.systemAlertRepo.FindAllByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		return nil, myerrors.WrapDomainError("systemAlertUsecase.GetAllAlerts", err)
	}

	return alerts, nil
}

func (u *systemAlertUsecase) MarkAlertAsRead(ctx context.Context, alertID int) error {
	err := u.systemAlertRepo.MarkAsRead(ctx, alertID)
	if err != nil {
		return myerrors.WrapDomainError("systemAlertUsecase.MarkAlertAsRead", err)
	}

	return nil
}

func (u *systemAlertUsecase) DeleteAlert(ctx context.Context, alertID int) error {
	err := u.systemAlertRepo.Delete(ctx, alertID)
	if err != nil {
		return myerrors.WrapDomainError("systemAlertUsecase.DeleteAlert", err)
	}

	return nil
}

func (u *systemAlertUsecase) GetAlertCount(ctx context.Context, projectID int) (int, error) {
	count, err := u.systemAlertRepo.CountActiveByProjectID(ctx, projectID)
	if err != nil {
		return 0, myerrors.WrapDomainError("systemAlertUsecase.GetAlertCount", err)
	}

	return count, nil
}

func (u *systemAlertUsecase) CheckAndCreateStorageAlert(ctx context.Context, projectID int, usagePercent float64) error {
	var severity string
	var title string
	var message string

	if usagePercent >= 95 {
		severity = SeverityCritical
		title = "ストレージ容量がほぼ満杯です"
		message = "ストレージ使用率が95%を超えています。すぐに容量を増やしてください。"
	} else if usagePercent >= 85 {
		severity = SeverityError
		title = "ストレージ容量が逼迫しています"
		message = "ストレージ使用率が85%を超えています。容量拡張を検討してください。"
	} else if usagePercent >= 75 {
		severity = SeverityWarning
		title = "ストレージ容量に注意が必要です"
		message = "ストレージ使用率が75%を超えています。使用状況を確認してください。"
	} else {
		return nil // アラート不要
	}

	metadata := map[string]interface{}{
		"usage_percent": usagePercent,
		"threshold":     "storage",
	}

	return u.CreateAlert(ctx, "storage", severity, title, message, projectID, metadata)
}

func (u *systemAlertUsecase) CheckAndCreateApiLimitAlert(ctx context.Context, projectID int, requestCount int, limit int) error {
	usagePercent := float64(requestCount) / float64(limit) * 100

	var severity string
	var title string
	var message string

	if usagePercent >= 95 {
		severity = SeverityCritical
		title = "APIリクエスト制限に達しています"
		message = "APIリクエスト数が制限の95%を超えています。制限解除を検討してください。"
	} else if usagePercent >= 85 {
		severity = SeverityError
		title = "APIリクエスト制限が近づいています"
		message = "APIリクエスト数が制限の85%を超えています。使用状況を確認してください。"
	} else if usagePercent >= 75 {
		severity = SeverityWarning
		title = "APIリクエスト制限に注意が必要です"
		message = "APIリクエスト数が制限の75%を超えています。制限拡大を検討してください。"
	} else {
		return nil // アラート不要
	}

	metadata := map[string]interface{}{
		"request_count": requestCount,
		"limit":         limit,
		"usage_percent": usagePercent,
		"threshold":     "api_limit",
	}

	return u.CreateAlert(ctx, "api_limit", severity, title, message, projectID, metadata)
}
