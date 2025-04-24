package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"w3st/domain/models"
	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/interfaces/services"
	"w3st/presenter"
	"w3st/usecase"
)

type UserController struct {
	userUsecase   usecase.UserUsecase
	authService   services.AuthService
	userPresenter presenter.UserPresenter
}

func NewUserController(userUsecase usecase.UserUsecase, authService services.AuthService, userPresenter presenter.UserPresenter) *UserController {
	return &UserController{
		userUsecase:   userUsecase,
		authService:   authService,
		userPresenter: userPresenter,
	}
}

func (c *UserController) Signup(ctx *gin.Context) {
	var input dto.SignupData

	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := &models.Users{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	// ユーザー登録
	newUser, err := c.userUsecase.Create(newUser, ctx.Request.Context())
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		// その他のエラー
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	// トークン生成
	token, err := c.authService.GenerateToken(newUser.ID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		// その他のエラー
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	// jwtトークンをクライアントに返す
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *UserController) Login(ctx *gin.Context) {
	var input dto.LoginData

	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ログイン
	user, err := c.userUsecase.FindByEmail(input.Email)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	//	token生成
	token, err := c.authService.GenerateToken(user.ID)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	// jwtトークンをクライアントに返す
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ユーザー情報の取得
	// userIDはstring型であることを確認
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.userUsecase.FindByID(userIDStr)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
