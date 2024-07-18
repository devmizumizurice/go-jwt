package services

import (
	"github.com/devmizumizurice/go-jwt/models/response"
	"github.com/devmizumizurice/go-jwt/repositories"
)

type UserServiceInterface interface {
	GetUserByID(id string) (*response.User, error)
	GetUserByEmail(email string) (*response.User, error)
}

type userService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface) UserServiceInterface {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUserByID(id string) (*response.User, error) {
	user, err := s.userRepository.FindByID(id)

	if err != nil {
		return nil, err
	}

	userResponse := &response.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResponse, nil
}

func (s *userService) GetUserByEmail(email string) (*response.User, error) {

	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		return nil, err
	}
	userResponse := &response.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResponse, nil
}
