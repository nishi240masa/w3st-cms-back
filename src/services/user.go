package services

import (
	"w3st/dto"
	"w3st/models"
	"w3st/repositories"

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

func Login(input dto.LoginData) (*models.User, error) {
    
        user, err := repositories.FindUser(input.Email)
        if err != nil {
            return nil, err
        }
    
        if user == nil {
            return nil, nil
        }
    
        // パスワードの照合
        err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
        if err != nil {
            return nil, err
        }
    
        return user, nil
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
