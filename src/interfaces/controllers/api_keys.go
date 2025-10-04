package controllers

import (
	"net/http"

	"w3st/dto"
	"w3st/usecase"

	"github.com/gin-gonic/gin"
)

type ApiKeyController struct {
	apiKeyUsecase usecase.ApiKeyUsecase
}

func NewApiKeyController(apiKeyUsecase usecase.ApiKeyUsecase) *ApiKeyController {
	return &ApiKeyController{
		apiKeyUsecase: apiKeyUsecase,
	}
}

func (ctrl *ApiKeyController) CreateApiKey(c *gin.Context) {
	var req dto.CreateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	projectID := c.GetInt("projectID")
	apiKey, err := ctrl.apiKeyUsecase.CreateApiKey(userID, projectID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"api_key": apiKey})
}