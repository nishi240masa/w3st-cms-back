package services

import (
	"os"
	"time"
	"w3st/dto"
	"w3st/models"
	"w3st/repositories"
	"w3st/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// func GetUsers() ([]models.User, error) {
// 	return repositorys.GetUsers()
// }

// func GetUser(userId string) (*models.User, error) {
// 	return repositorys.GetUser(userId)
// }


func Signup(input dto.SignupData) (*models.User, error) {

    // passwordのハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // 新しいユーザーを作成
    newUser := &models.User{
        Name:    input.Name,
        Email:   input.Email,
        Password: string(hashedPassword),
    }

    err = repositories.CreateUser(newUser)

    if err != nil {
        return nil, err
    }

    return newUser, nil
}

func Login(input dto.LoginData) (*string, error) {
    
        foundUser, err := repositories.FindUser(input.Email)
        if err != nil {
            return nil, err
        }
    
        if foundUser == nil {
            return nil, nil
        }
    
        // パスワードの照合
        err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(input.Password))
        if err != nil {
            return nil, err
        }

        // tokenの作成

        userID, err := utils.UuidToUint(foundUser.ID)
        if err != nil {
            return nil, err
        }

        token, err := CreateToken(userID, foundUser.Email)
        if err != nil {
            return nil, err
        }

    
        return token, nil
    }


func CreateToken(userId uint, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}



// func CreateUser(input dto.CreateUserData) (*models.User, error) {

//     // 既に登録されているか確認
//     user, err := repositorys.GetUserByName(input.Name)
//     if err != nil {
//         return nil, err
//     }

//     if user != nil {
//         // ユーザーが既に存在する場合の処理
//         return nil, errors.New("user already exists")
//     }

//     // 新しいユーザーを作成
//     newUser := &models.User{
//         Name:    input.Name,
//         Email:   input.Email,
//         Password: input.Password,
//     }

// 	err = repositorys.CreateUser(newUser)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return newUser, nil
// }
