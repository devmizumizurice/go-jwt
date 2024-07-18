package repositories

import (
	"github.com/devmizumizurice/go-jwt/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(user *models.User) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	res := r.db.Create(user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	res := r.db.Where("id = ?", id).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	res := r.db.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
