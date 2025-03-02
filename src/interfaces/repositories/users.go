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
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}
	
	return &user, nil
	
}
