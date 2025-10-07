package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"w3st/dto"
	"w3st/usecase"
)

type ApiKeyController struct {
	apiKeyUsecase usecase.ApiKeyUsecase
}

func NewApiKeyController(apiKeyUsecase usecase.ApiKeyUsecase) *ApiKeyController {
	return &ApiKeyController{
		apiKeyUsecase: apiKeyUsecase,
	}
}

func (c *ApiKeyController) CreateApiKey(ctx *gin.Context) {
	var req dto.CreateApiKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	projectID, exists := ctx.Get("projectID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Project ID not found"})
		return
	}

	projectIDInt, ok := projectID.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid project ID"})
		return
	}

	apiKey, err := c.apiKeyUsecase.CreateApiKey(userIDStr, projectIDInt, req.Name, req.CollectionIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"api_key": apiKey})
}

func (c *ApiKeyController) ValidateApiKey(ctx *gin.Context) {
	apiKey := ctx.Query("key")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is required"})
		return
	}

	token, err := c.apiKeyUsecase.ValidateApiKey(apiKey)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
