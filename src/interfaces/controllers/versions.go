package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
)

type VersionController struct {
	BaseController
	versionUsecase usecase.VersionUsecase
}

func NewVersionController(versionUsecase usecase.VersionUsecase) *VersionController {
	return &VersionController{
		versionUsecase: versionUsecase,
	}
}

func (c *VersionController) CreateVersion(ctx *gin.Context) {
	var input dto.CreateVersion

	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	contentUUID, err := uuid.Parse(input.ContentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID format"})
		return
	}

	// JSONデータをパース
	var data interface{}
	if err := json.Unmarshal([]byte(input.Data), &data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	// バージョン作成
	version, err := c.versionUsecase.CreateVersion(ctx.Request.Context(), userUUID, contentUUID, data)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// JSONデータを文字列化
	dataBytes, err := json.Marshal(version.Data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	dataStr := string(dataBytes)

	// レスポンス
	response := dto.VersionResponse{
		ID:        version.ID.String(),
		ContentID: version.ContentID.String(),
		Version:   version.Version,
		Data:      dataStr,
		UserID:    version.UserID.String(),
		CreatedAt: version.CreatedAt.String(),
		UpdatedAt: version.UpdatedAt.String(),
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *VersionController) GetVersionsByContentID(ctx *gin.Context) {
	contentID := ctx.Param("contentID")

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	contentUUID, err := uuid.Parse(contentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID format"})
		return
	}

	// バージョン一覧取得
	versions, err := c.versionUsecase.GetVersionsByContentID(ctx.Request.Context(), userUUID, contentUUID)
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
	responses := make([]dto.VersionResponse, 0, len(versions))
	for _, version := range versions {
		dataBytes, err := json.Marshal(version.Data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
			return
		}
		dataStr := string(dataBytes)
		responses = append(responses, dto.VersionResponse{
			ID:        version.ID.String(),
			ContentID: version.ContentID.String(),
			Version:   version.Version,
			Data:      dataStr,
			UserID:    version.UserID.String(),
			CreatedAt: version.CreatedAt.String(),
			UpdatedAt: version.UpdatedAt.String(),
		})
	}
	ctx.JSON(http.StatusOK, responses)
}

func (c *VersionController) GetLatestVersion(ctx *gin.Context) {
	contentID := ctx.Param("contentID")

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	contentUUID, err := uuid.Parse(contentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID format"})
		return
	}

	// 最新バージョン取得
	version, err := c.versionUsecase.GetLatestVersion(ctx.Request.Context(), userUUID, contentUUID)
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
	dataBytes, err := json.Marshal(version.Data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	dataStr := string(dataBytes)
	response := dto.VersionResponse{
		ID:        version.ID.String(),
		ContentID: version.ContentID.String(),
		Version:   version.Version,
		Data:      dataStr,
		UserID:    version.UserID.String(),
		CreatedAt: version.CreatedAt.String(),
		UpdatedAt: version.UpdatedAt.String(),
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *VersionController) RestoreVersion(ctx *gin.Context) {
	contentID := ctx.Param("contentID")
	versionID := ctx.Param("versionID")

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	contentUUID, err := uuid.Parse(contentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID format"})
		return
	}

	versionUUID, err := uuid.Parse(versionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID format"})
		return
	}

	// バージョン復元
	version, err := c.versionUsecase.RestoreVersion(ctx.Request.Context(), userUUID, contentUUID, versionUUID)
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
	dataBytes, err := json.Marshal(version.Data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	dataStr := string(dataBytes)
	response := dto.VersionResponse{
		ID:        version.ID.String(),
		ContentID: version.ContentID.String(),
		Version:   version.Version,
		Data:      dataStr,
		UserID:    version.UserID.String(),
		CreatedAt: version.CreatedAt.String(),
		UpdatedAt: version.UpdatedAt.String(),
	}
	ctx.JSON(http.StatusOK, response)
}
