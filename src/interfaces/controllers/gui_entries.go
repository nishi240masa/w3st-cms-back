package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"w3st/domain/models"
	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
)

type GUIEntriesController struct {
	entriesUsecase usecase.EntriesUsecase
}

func NewGUIEntriesController(entriesUsecase usecase.EntriesUsecase) *GUIEntriesController {
	return &GUIEntriesController{
		entriesUsecase: entriesUsecase,
	}
}

func (c *GUIEntriesController) GetEntries(ctx *gin.Context) {
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

	// entriesを取得
	entries, err := c.entriesUsecase.GetEntriesByCollectionId(collectionIdInt, projectID)
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

func (c *GUIEntriesController) CreateEntry(ctx *gin.Context) {
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

	// リクエストのバインド
	var input dto.CreateEntry
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// dataをJSONに変換
	dataBytes, err := json.Marshal(input.Data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	// entryを作成
	newEntry := &models.Entry{
		ProjectID:    projectID,
		CollectionID: collectionIdInt,
		Data:         string(dataBytes),
	}

	// entryを作成
	err = c.entriesUsecase.CreateEntry(newEntry, projectID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Entry created successfully"})
}
