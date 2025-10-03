package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
)

type AuditController struct {
	auditUsecase usecase.AuditUsecase
}

func NewAuditController(auditUsecase usecase.AuditUsecase) *AuditController {
	return &AuditController{
		auditUsecase: auditUsecase,
	}
}

// getUserUUID extracts and parses userID from gin context
func (c *AuditController) getUserUUID(ctx *gin.Context) uuid.UUID {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return uuid.Nil
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return uuid.Nil
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return uuid.Nil
	}
	return userUUID
}

func (c *AuditController) LogAction(ctx *gin.Context) {
	var input dto.CreateAuditLog

	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// アクションログ
	err := c.auditUsecase.LogAction(ctx.Request.Context(), userUUID, input.Action, input.Resource, input.Details)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Action logged"})
}

func (c *AuditController) GetLogsByUser(ctx *gin.Context) {
	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// ユーザーのログ取得
	logs, err := c.auditUsecase.GetLogsByUser(ctx.Request.Context(), userUUID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス
	responses := make([]dto.AuditLogResponse, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, dto.AuditLogResponse{
			ID:        log.ID.String(),
			UserID:    log.UserID.String(),
			Action:    log.Action,
			Resource:  log.Resource,
			Details:   log.Details,
			CreatedAt: log.CreatedAt.String(),
		})
	}
	ctx.JSON(http.StatusOK, responses)
}

func (c *AuditController) GetLogsByAction(ctx *gin.Context) {
	action := ctx.Param("action")

	// アクション別のログ取得
	logs, err := c.auditUsecase.GetLogsByAction(ctx.Request.Context(), action)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス
	responses := make([]dto.AuditLogResponse, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, dto.AuditLogResponse{
			ID:        log.ID.String(),
			UserID:    log.UserID.String(),
			Action:    log.Action,
			Resource:  log.Resource,
			Details:   log.Details,
			CreatedAt: log.CreatedAt.String(),
		})
	}
	ctx.JSON(http.StatusOK, responses)
}

func (c *AuditController) GetAllLogs(ctx *gin.Context) {
	// 全ログ取得
	logs, err := c.auditUsecase.GetAllLogs(ctx.Request.Context())
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス
	responses := make([]dto.AuditLogResponse, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, dto.AuditLogResponse{
			ID:        log.ID.String(),
			UserID:    log.UserID.String(),
			Action:    log.Action,
			Resource:  log.Resource,
			Details:   log.Details,
			CreatedAt: log.CreatedAt.String(),
		})
	}
	ctx.JSON(http.StatusOK, responses)
}
