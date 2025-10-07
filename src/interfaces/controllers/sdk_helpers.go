package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseCollectionRequest parses collectionId param, projectID and collectionIds from context.
// Returns status != 0 and errMsg when an HTTP error should be returned by the caller.
func parseCollectionRequest(ctx *gin.Context) (collectionIdInt int, projectID int, collectionIds []int, status int, errMsg string) {
	collectionId := ctx.Param("collectionId")
	collectionIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		return 0, 0, nil, http.StatusBadRequest, "Invalid Collection ID"
	}

	projectID = ctx.GetInt("projectID")

	collectionIdsInterface, exists := ctx.Get("collectionIds")
	if !exists {
		return 0, projectID, nil, http.StatusUnauthorized, "Collection IDs not found in context"
	}
	ids, ok := collectionIdsInterface.([]int)
	if !ok {
		return 0, projectID, nil, http.StatusInternalServerError, "Invalid collection IDs format"
	}
	return collectionIdInt, projectID, ids, 0, ""
}
