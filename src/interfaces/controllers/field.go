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
)

type FieldController struct {
	fieldUsecase usecase.FieldUsecase
}

func NewFieldController(fieldUsecase usecase.FieldUsecase) *FieldController {
	return &FieldController{
		fieldUsecase: fieldUsecase,
	}
}

func (f *FieldController) Create(ctx *gin.Context) {
	// パラメータから取得
	collectionId := ctx.Param("collectionId")
	if collectionId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Collection ID is required"})
		return
	}
	// int型に変換
	collectionIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collection ID"})
		return
	}

	// projectID
	projectID, exists := ctx.Get("projectID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Project ID not found in context"})
		return
	}

	// projectIDはint型であることを確認
	projectIDInt, ok := projectID.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// リクエストのバインド
	var input dto.CreateField
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// フィールドの作成
	newField := &models.FieldData{
		ProjectID:    projectIDInt,
		CollectionID: collectionIdInt,
		FieldID:      input.FieldID,
		ViewName:     input.ViewName,
		FieldType:    input.FieldType,
		IsRequired:   input.IsRequired,
		DefaultValue: input.DefaultValue,
	}

	err = f.fieldUsecase.Create(projectIDInt, newField)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンスを返
	ctx.JSON(http.StatusOK, gin.H{"message": "Field created successfully"})
}

func (f *FieldController) Update(ctx *gin.Context) {
	// フィールドIDを取得
	fieldId := ctx.Param("fieldId")
	if fieldId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Field ID is required"})
		return
	}

	//　コレクションIDを取得
	collectionId := ctx.Param("collectionId")
	if collectionId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Collection ID is required"})
		return
	}

	// int型に変換
	collectionIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collection ID"})
		return
	}

	// projectID
	projectID, exists := ctx.Get("projectID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Project ID not found in context"})
		return
	}

	// projectIDはint型であることを確認
	projectIDInt, ok := projectID.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// リクエストのバインド
	var input dto.UpdateField
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newField := &models.FieldData{
		ProjectID:    projectIDInt,
		CollectionID: collectionIdInt,
		FieldID:      input.FieldID,
		ViewName:     input.ViewName,
		FieldType:    input.FieldType,
		IsRequired:   input.IsRequired,
		DefaultValue: input.DefaultValue,
	}

	// フィールドの更新
	err = f.fieldUsecase.Update(projectIDInt, newField)
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

func (f *FieldController) Delete(ctx *gin.Context) {
	// パラメータから取得
	collectionId := ctx.Param("collectionId")
	if collectionId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Collection ID is required"})
		return
	}
	fieldId := ctx.Param("fieldId")
	if fieldId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Field ID is required"})
		return
	}

	// projectID
	projectID, exists := ctx.Get("projectID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Project ID not found in context"})
		return
	}

	projectIDInt, ok := projectID.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// フィールドを削除
	err := f.fieldUsecase.Delete(projectIDInt, fieldId)
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
