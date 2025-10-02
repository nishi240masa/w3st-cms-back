package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"w3st/domain/models"
	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
	"w3st/utils"
)

type CollectionsController struct {
	collectionUsecase usecase.CollectionsUsecase
}

func NewCollectionsController(collectionUsecase usecase.CollectionsUsecase) *CollectionsController {
	return &CollectionsController{
		collectionUsecase: collectionUsecase,
	}
}

func (c *CollectionsController) MakeCollection(ctx *gin.Context) {
	// リクエストのバインド
	var input dto.MakeCollection
	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// userID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// userIDはstring型であることを確認
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userUuid, err := utils.StringToUUID(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// コレクション
	newCollection := &models.ApiCollection{
		UserID:      userUuid,
		Name:        input.Name,
		Description: input.Description,
	}

	//	collectionを作成
	err = c.collectionUsecase.Make(newCollection)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	// レスポンスを返す
	ctx.JSON(http.StatusOK, gin.H{"message": "Collection created successfully"})
}

func (c *CollectionsController) GetCollectionByUserId(ctx *gin.Context) {
	// ユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// userIDはstring型であることを確認
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userUuid, err := utils.StringToUUID(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// コレクションを取得
	collection, err := c.collectionUsecase.GetCollectionByUserId(userUuid)
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

func (c *CollectionsController) GetCollectionsByCollectionId(ctx *gin.Context) {
	// コレクションIDを取得
	collectionId := ctx.Param("collectionId")

	//　int型に変換
	collectionIdInt, err := strconv.Atoi(collectionId)

	// userIdを取得
	// ユーザーIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// userIDはstring型であることを確認
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userUuid, err := utils.StringToUUID(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// コレクションを取得
	collection, err := c.collectionUsecase.GetCollectionsByCollectionId(collectionIdInt, userUuid)
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
