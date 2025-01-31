package repositories

import "w3st/models"




func CreateUser(newUser *models.User) error {
	result := db.Create(newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}


func FindUser (email string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}