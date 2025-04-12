package controllers

import (
	"net/http"
	"w3st/domain/models"
	"w3st/dto"
	"w3st/interfaces/services"
	"w3st/presenter"
	"w3st/usecase"

	"github.com/gin-gonic/gin"
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

	//fmt.Println("input:", input)
	//fmt.Println("newUser:", newUser.Name)

	// ユーザー登録
	token, err := c.userUsecase.Create(newUser, ctx.Request.Context())
	if err != nil {
		err := ErrorHandle(err)
		ctx.JSON(HttpStatusCodeFromConnectCode(err.Code()), gin.H{"error": err.Error()})
		return
	}

	// jwtトークンをクライアントに返す
	ctx.JSON(http.StatusOK, token)
}

func (c *UserController) Login(ctx *gin.Context) {
	var input dto.LoginData

	// リクエストのバインド
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ログイン
	token, err := c.userUsecase.FindByEmail(input.Email)
	if err != nil {
		err := ErrorHandle(err)
		ctx.JSON(HttpStatusCodeFromConnectCode(err.Code()), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})

}
