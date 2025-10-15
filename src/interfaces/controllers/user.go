package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"w3st/domain/models"
	"w3st/dto"
	myerrors "w3st/errors"
	"w3st/presenter"
	"w3st/usecase"
)

type UserController struct {
	userUsecase    usecase.UserUsecase
	jwtAuthUsecase usecase.JwtUsecase
	userPresenter  presenter.UserPresenter
}

func NewUserController(userUsecase usecase.UserUsecase, jwtAuthUsecase usecase.JwtUsecase, userPresenter presenter.UserPresenter) *UserController {
	return &UserController{
		userUsecase:    userUsecase,
		jwtAuthUsecase: jwtAuthUsecase,
		userPresenter:  userPresenter,
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
	token, err := c.jwtAuthUsecase.GenerateToken(newUser.ID)
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
	token, err := c.jwtAuthUsecase.GenerateToken(user.ID)
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

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userUsecase.GetAllUsers()
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (c *UserController) UpdateUserById(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var input dto.UpdateUserData
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userUsecase.FindByID(userId)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Update user fields
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	err = c.userUsecase.Update(user, ctx.Request.Context())
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

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	err := c.userUsecase.DeleteUser(userId)
	if err != nil {
		var domainErr *myerrors.DomainError
		if errors.As(err, &domainErr) {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input dto.UpdateUserData
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Update user fields
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	err = c.userUsecase.Update(user, ctx.Request.Context())
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
