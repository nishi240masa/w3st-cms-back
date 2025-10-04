package controllers

import (
	"errors"
	"net/http"
	"strconv"

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

	// コレクションを取得
	collection, err := c.collectionUsecase.GetCollectionByProjectId(projectID)
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
	// コレクションIDを取得
	collectionId := ctx.Param("collectionId")

	//　int型に変換
	collectionIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collection ID"})
		return
	}

	// プロジェクトIDを取得
	projectID := ctx.GetInt("projectID")

	// コレクションを取得
	collection, err := c.collectionUsecase.GetCollectionsByCollectionId(collectionIdInt, projectID)
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
