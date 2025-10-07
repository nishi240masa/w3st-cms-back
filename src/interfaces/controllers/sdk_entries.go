package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	myerrors "w3st/errors"
	"w3st/usecase"
)

type SDKEntriesController struct {
	entriesUsecase usecase.EntriesUsecase
}

func NewSDKEntriesController(entriesUsecase usecase.EntriesUsecase) *SDKEntriesController {
	return &SDKEntriesController{
		entriesUsecase: entriesUsecase,
	}
}

func (c *SDKEntriesController) GetEntries(ctx *gin.Context) {
	// collectionIdを取得
	collectionId := ctx.Param("collectionId")

	// int型に変換
	collectionIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collection ID"})
		return
	}

	// プロジェクトIDを取得
	projectID := ctx.GetInt("projectID")

	// collectionIdsを取得
	collectionIdsInterface, exists := ctx.Get("collectionIds")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	collectionIds, ok := collectionIdsInterface.([]int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// entriesを取得
	entries, err := c.entriesUsecase.GetEntriesByCollectionIdForSDK(collectionIdInt, projectID, collectionIds)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
