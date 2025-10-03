package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/usecase"
)

type MediaController struct {
	mediaUsecase usecase.MediaUsecase
}

func NewMediaController(mediaUsecase usecase.MediaUsecase) *MediaController {
	return &MediaController{
		mediaUsecase: mediaUsecase,
	}
}

// getUserUUID extracts and parses userID from gin context
func (c *MediaController) getUserUUID(ctx *gin.Context) uuid.UUID {
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

func (c *MediaController) Upload(ctx *gin.Context) {
	var input dto.CreateMedia

	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// メディアアップロード
	media, err := c.mediaUsecase.Upload(ctx.Request.Context(), userUUID, input.Name, input.Type, input.Path, input.Size)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス
	response := dto.MediaResponse{
		ID:        media.ID.String(),
		Name:      media.Name,
		Type:      media.Type,
		Path:      media.Path,
		Size:      media.Size,
		UserID:    media.UserID.String(),
		CreatedAt: media.CreatedAt.String(),
		UpdatedAt: media.UpdatedAt.String(),
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *MediaController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// メディア取得
	media, err := c.mediaUsecase.GetByID(ctx.Request.Context(), userUUID, id)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス
	response := dto.MediaResponse{
		ID:        media.ID.String(),
		Name:      media.Name,
		Type:      media.Type,
		Path:      media.Path,
		Size:      media.Size,
		UserID:    media.UserID.String(),
		CreatedAt: media.CreatedAt.String(),
		UpdatedAt: media.UpdatedAt.String(),
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *MediaController) GetByUserID(ctx *gin.Context) {
	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// ユーザーのメディア一覧取得
	medias, err := c.mediaUsecase.GetByUserID(ctx.Request.Context(), userUUID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// レスポンス
	responses := make([]dto.MediaResponse, 0, len(medias))
	for _, media := range medias {
		responses = append(responses, dto.MediaResponse{
			ID:        media.ID.String(),
			Name:      media.Name,
			Type:      media.Type,
			Path:      media.Path,
			Size:      media.Size,
			UserID:    media.UserID.String(),
			CreatedAt: media.CreatedAt.String(),
			UpdatedAt: media.UpdatedAt.String(),
		})
	}
	ctx.JSON(http.StatusOK, responses)
}

func (c *MediaController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	userUUID := c.getUserUUID(ctx)
	if userUUID == uuid.Nil {
		return
	}

	// メディア削除
	err := c.mediaUsecase.Delete(ctx.Request.Context(), userUUID, id)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
