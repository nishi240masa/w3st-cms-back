package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	myerrors "w3st/errors"
)

type BaseController struct{}

func (c *BaseController) getUserUUID(ctx *gin.Context) uuid.UUID {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return uuid.Nil
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return uuid.Nil
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return uuid.Nil
	}
	return userUUID
}

func (c *BaseController) marshalDataToString(data interface{}) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", myerrors.NewDomainError(myerrors.ErrorUnknown, err)
	}
	return string(dataBytes), nil
}
