package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	myerrors "w3st/errors"
	"w3st/usecase"
)

type SDKCollectionsController struct {
	collectionUsecase usecase.CollectionsUsecase
}

func NewSDKCollectionsController(collectionUsecase usecase.CollectionsUsecase) *SDKCollectionsController {
	return &SDKCollectionsController{
		collectionUsecase: collectionUsecase,
	}
}

// GetCollectionByProjectId - SDK用：プロジェクトのコレクション一覧取得
func (c *SDKCollectionsController) GetCollectionByProjectId(ctx *gin.Context) {
	// プロジェクトIDを取得
	projectID := ctx.GetInt("projectID")
	// 許可されたコレクションIDを取得
	collectionIds, exists := ctx.Get("collectionIds")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Collection IDs not found in context"})
		return
	}
	collectionIdsSlice, ok := collectionIds.([]int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid collection IDs format"})
		return
	}

	// コレクションを取得（フィルタリング済み）
	collection, err := c.collectionUsecase.GetCollectionByProjectIdForSDK(projectID, collectionIdsSlice)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, collection)
}

// GetCollectionsByCollectionId - SDK用：コレクション詳細取得
func (c *SDKCollectionsController) GetCollectionsByCollectionId(ctx *gin.Context) {
	collectionIdInt, projectID, collectionIds, status, errMsg := parseCollectionRequest(ctx)
	if status != 0 {
		ctx.JSON(status, gin.H{"error": errMsg})
		return
	}

	collection, err := c.collectionUsecase.GetCollectionsByCollectionIdForSDK(collectionIdInt, projectID, collectionIds)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, collection)
}
