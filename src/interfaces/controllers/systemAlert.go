package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
)

type SystemAlertController struct {
	systemAlertUsecase usecase.SystemAlertUsecase
}

func NewSystemAlertController(systemAlertUsecase usecase.SystemAlertUsecase) *SystemAlertController {
	return &SystemAlertController{
		systemAlertUsecase: systemAlertUsecase,
	}
}

func (c *SystemAlertController) GetAlerts(ctx *gin.Context) {
	// プロジェクトID取得
	projectID := ctx.GetInt("projectID")

	// クエリパラメータ取得
	var query dto.GetSystemAlertsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// アラート取得
	alerts, err := c.systemAlertUsecase.GetAllAlerts(ctx.Request.Context(), projectID, query.Limit, query.Offset)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス変換
	response := make([]dto.SystemAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		response = append(response, dto.SystemAlertResponse{
			ID:        alert.ID,
			AlertType: alert.AlertType,
			Severity:  alert.Severity,
			Title:     alert.Title,
			Message:   alert.Message,
			ProjectID: alert.ProjectID,
			IsActive:  alert.IsActive,
			IsRead:    alert.IsRead,
			Metadata:  alert.Metadata,
			CreatedAt: alert.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: alert.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *SystemAlertController) GetActiveAlerts(ctx *gin.Context) {
	// プロジェクトID取得
	projectID := ctx.GetInt("projectID")

	// アクティブなアラート取得
	alerts, err := c.systemAlertUsecase.GetActiveAlerts(ctx.Request.Context(), projectID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス変換
	response := make([]dto.SystemAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		response = append(response, dto.SystemAlertResponse{
			ID:        alert.ID,
			AlertType: alert.AlertType,
			Severity:  alert.Severity,
			Title:     alert.Title,
			Message:   alert.Message,
			ProjectID: alert.ProjectID,
			IsActive:  alert.IsActive,
			IsRead:    alert.IsRead,
			Metadata:  alert.Metadata,
			CreatedAt: alert.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: alert.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *SystemAlertController) CreateAlert(ctx *gin.Context) {
	// リクエストバインド
	var input dto.CreateSystemAlertRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID := ctx.GetInt("projectID")

	// アラート作成
	err := c.systemAlertUsecase.CreateAlert(ctx.Request.Context(), input.AlertType, input.Severity, input.Title, input.Message, projectID, input.Metadata)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Alert created successfully"})
}

func (c *SystemAlertController) MarkAsRead(ctx *gin.Context) {
	alertID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Alert ID"})
		return
	}

	err = c.systemAlertUsecase.MarkAlertAsRead(ctx.Request.Context(), alertID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Alert marked as read"})
}

func (c *SystemAlertController) DeleteAlert(ctx *gin.Context) {
	alertID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Alert ID"})
		return
	}

	err = c.systemAlertUsecase.DeleteAlert(ctx.Request.Context(), alertID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})
}

func (c *SystemAlertController) GetAlertCount(ctx *gin.Context) {
	projectID := ctx.GetInt("projectID")

	count, err := c.systemAlertUsecase.GetAlertCount(ctx.Request.Context(), projectID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SystemAlertCountResponse{Count: count})
}
