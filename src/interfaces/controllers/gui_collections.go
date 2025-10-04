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

type GUICollectionsController struct {
	collectionUsecase usecase.CollectionsUsecase
	fieldUsecase      usecase.FieldUsecase
}

func NewGUICollectionsController(collectionUsecase usecase.CollectionsUsecase, fieldUsecase usecase.FieldUsecase) *GUICollectionsController {
	return &GUICollectionsController{
		collectionUsecase: collectionUsecase,
		fieldUsecase:      fieldUsecase,
	}
}

// MakeCollection - GUI用：コレクション作成
func (c *GUICollectionsController) MakeCollection(ctx *gin.Context) {
	// リクエストのバインド
	var input dto.MakeCollection
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

	// プロジェクトIDを取得
	projectID := ctx.GetInt("projectID")

	// コレクション
	newCollection := &models.ApiCollection{
		UserID:      userUuid,
		ProjectID:   projectID,
		Name:        input.Name,
		Description: input.Description,
	}

	// collectionを作成
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

// CreateField - GUI用：フィールド作成
func (c *GUICollectionsController) CreateField(ctx *gin.Context) {
	collectionId := ctx.Param("collectionId")

	// int型に変換
	collectionIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collection ID"})
		return
	}

	var input dto.CreateField
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// プロジェクトIDを取得
	projectID := ctx.GetInt("projectID")

	// DTOからモデルに変換
	fieldData := &models.FieldData{
		ProjectID:    projectID,
		CollectionID: collectionIdInt,
		FieldID:      input.FieldID,
		ViewName:     input.ViewName,
		FieldType:    input.FieldType,
		IsRequired:   input.IsRequired,
		DefaultValue: input.DefaultValue,
	}

	err = c.fieldUsecase.Create(projectID, fieldData)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Field created successfully"})
}

// UpdateField - GUI用：フィールド更新
func (c *GUICollectionsController) UpdateField(ctx *gin.Context) {
	collectionId := ctx.Param("collectionId")
	fieldId := ctx.Param("fieldId")

	// int型に変換
	collectionIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collection ID"})
		return
	}

	fieldIdInt, err := strconv.Atoi(fieldId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Field ID"})
		return
	}

	var input dto.UpdateField
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// プロジェクトIDを取得
	projectID := ctx.GetInt("projectID")

	// DTOからモデルに変換
	fieldData := &models.FieldData{
		ID:           fieldIdInt,
		ProjectID:    projectID,
		CollectionID: collectionIdInt,
		FieldID:      input.FieldID,
		ViewName:     input.ViewName,
		FieldType:    input.FieldType,
		IsRequired:   input.IsRequired,
		DefaultValue: input.DefaultValue,
	}

	err = c.fieldUsecase.Update(projectID, fieldData)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Field updated successfully"})
}

// DeleteField - GUI用：フィールド削除
func (c *GUICollectionsController) DeleteField(ctx *gin.Context) {
	fieldId := ctx.Param("fieldId")

	// プロジェクトIDを取得
	projectID := ctx.GetInt("projectID")

	err := c.fieldUsecase.Delete(projectID, fieldId)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Field deleted successfully"})
}
