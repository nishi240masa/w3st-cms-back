package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
)

type AuditController struct {
	BaseController
	auditUsecase usecase.AuditUsecase
}

func NewAuditController(auditUsecase usecase.AuditUsecase) *AuditController {
	return &AuditController{
		auditUsecase: auditUsecase,
	}
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

	// プロジェクトIDを取得（APIキー認証の場合はコンテキストから、GUIの場合はデフォルト）
	projectID := 1 // デフォルト
	if pID, exists := ctx.Get("projectID"); exists {
		if pid, ok := pID.(int); ok {
			projectID = pid
		}
	}

	// アクションログ
	err := c.auditUsecase.LogActionWithProject(ctx.Request.Context(), userUUID, projectID, input.Action, input.Resource, input.Details)
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

func (c *AuditController) GetLogsByProject(ctx *gin.Context) {
	projectIDStr := ctx.Param("projectId")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// クエリパラメータからlimitとoffsetを取得、デフォルト値設定
	limitStr := ctx.DefaultQuery("limit", "50")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// プロジェクトのログ取得
	logs, err := c.auditUsecase.GetLogsByProjectWithLimit(ctx.Request.Context(), projectID, limit, offset)
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
	// クエリパラメータからlimitとoffsetを取得、デフォルト値設定
	limitStr := ctx.DefaultQuery("limit", "50")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// 全ログ取得（管理者権限が必要）
	logs, err := c.auditUsecase.GetAllLogs(ctx.Request.Context(), limit, offset)
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
