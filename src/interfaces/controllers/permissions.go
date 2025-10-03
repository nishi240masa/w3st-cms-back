package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
)

type PermissionActionOptions struct {
	Action         func(context.Context, uuid.UUID, string, string) error
	SuccessMessage string
	StatusCode     int
}

type PermissionController struct {
	BaseController
	permissionUsecase usecase.PermissionUsecase
}

func NewPermissionController(permissionUsecase usecase.PermissionUsecase) *PermissionController {
	return &PermissionController{
		permissionUsecase: permissionUsecase,
	}
}

// handlePermissionAction handles grant or revoke permission actions
func (c *PermissionController) handlePermissionAction(ctx *gin.Context, opts PermissionActionOptions) {
	var input dto.CreatePermission

	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// 権限操作
	err := opts.Action(ctx.Request.Context(), userUUID, input.Permission, input.Resource)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(opts.StatusCode, gin.H{"message": opts.SuccessMessage})
}

func (c *PermissionController) CheckPermission(ctx *gin.Context) {
	permission := ctx.Query("permission")
	resource := ctx.Query("resource")

	if permission == "" || resource == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "permission and resource are required"})
		return
	}

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// 権限チェック
	hasPermission, err := c.permissionUsecase.CheckPermission(ctx.Request.Context(), userUUID, permission, resource)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"has_permission": hasPermission})
}

func (c *PermissionController) GrantPermission(ctx *gin.Context) {
	c.handlePermissionAction(ctx, PermissionActionOptions{
		Action:         c.permissionUsecase.GrantPermission,
		SuccessMessage: "Permission granted",
		StatusCode:     http.StatusCreated,
	})
}

func (c *PermissionController) RevokePermission(ctx *gin.Context) {
	c.handlePermissionAction(ctx, PermissionActionOptions{
		Action:         c.permissionUsecase.RevokePermission,
		SuccessMessage: "Permission revoked",
		StatusCode:     http.StatusOK,
	})
}

func (c *PermissionController) GetPermissionsByUser(ctx *gin.Context) {
	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// ユーザー権限一覧取得
	permissions, err := c.permissionUsecase.GetPermissionsByUser(ctx.Request.Context(), userUUID)
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
	responses := make([]dto.PermissionResponse, 0, len(permissions))
	for _, perm := range permissions {
		responses = append(responses, dto.PermissionResponse{
			ID:         perm.ID.String(),
			UserID:     perm.UserID.String(),
			Permission: perm.Permission,
			Resource:   perm.Resource,
			CreatedAt:  perm.CreatedAt.String(),
			UpdatedAt:  perm.UpdatedAt.String(),
		})
	}
	ctx.JSON(http.StatusOK, responses)
}
