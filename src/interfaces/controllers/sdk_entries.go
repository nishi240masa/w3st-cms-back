package controllers

import (
	"errors"
	"net/http"

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
	collectionIdInt, projectID, collectionIds, status, errMsg := parseCollectionRequest(ctx)
	if status != 0 {
		ctx.JSON(status, gin.H{"error": errMsg})
		return
	}

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
