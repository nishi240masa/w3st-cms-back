package repositories

import (
	"w3st/models"

	"gorm.io/gorm"
)


type UserRepository interface {
	Create(newUser *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userRepository {
	return userRepository{db}
}


// Create
func (r *userRepository) Create(newUser *models.User) error {
	return r.db.Create(newUser).Error
}

// Find
func (r *userRepository) FindByEmail (email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
	
}
