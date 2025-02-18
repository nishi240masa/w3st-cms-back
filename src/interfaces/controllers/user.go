package controllers

import (
	"net/http"
	"w3st/dto"
	"w3st/services"

	"github.com/gin-gonic/gin"
)

// func extractJWTFromHeader(c *gin.Context) (string, error) {
// 	authHeader := c.GetHeader("Authorization")
// 	if authHeader == "" {
// 		return "", errors.New("missing Authorization header")
// 	}

// 	parts := strings.Split(authHeader, " ")
// 	if len(parts) != 2 || parts[0] != "Bearer" {
// 		return "", errors.New("invalid Authorization header format")
// 	}

// 	return parts[1], nil
// }

// func GetUsers(c *gin.Context) {

// 	users, err := services.GetUsers()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, users)
// }

// func GetUser(c *gin.Context) {

// 	userId := c.Param("id")

// 	user, err := services.GetUser(userId)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

func Signup(c *gin.Context) {
	var input dto.SignupData

	// リクエストのバインド
	if  err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザー登録
	user, err := services.Signup(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context) {
	var input dto.LoginData

	// リクエストのバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ログイン
	token, err := services.Login(input)
		if err != nil {
			if err.Error() == "User not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
		}


	c.JSON(http.StatusOK, token)
}

// func CreateUser(c *gin.Context) {

// 	var input dto.CreateUserData
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user, err := services.CreateUser(input)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }
