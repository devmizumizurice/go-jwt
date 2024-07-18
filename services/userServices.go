package services

import (
	"errors"

	"github.com/devmizumizurice/go-jwt/models"
	"github.com/devmizumizurice/go-jwt/models/response"
	"github.com/devmizumizurice/go-jwt/repositories"
	"github.com/devmizumizurice/go-jwt/utils"
)

type UserServiceInterface interface {
	CreateUser(user *models.User) (*response.User, error)
	VerifyUserEmail(email string, password string) (*response.Token, error)
	GetUserByID(id string) (*response.User, error)
	GetUserByEmail(email string) (*response.User, error)
}

type userService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface) UserServiceInterface {
	return &userService{userRepository: userRepository}
}

func (s *userService) CreateUser(user *models.User) (*response.User, error) {
	existingUser, _ := s.userRepository.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("EMAIL_ALREADY_EXISTS")
	}

	hash, err := utils.PasswordEncrypt(user.Password)

	if err != nil {
		return nil, errors.New("ERROR_ENCRYPTING_PASSWORD")
	}

	user.Password = hash

	user, err = s.userRepository.Create(user)

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

func (s *userService) VerifyUserEmail(email string, password string) (*response.Token, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("INVALID_EMAIL_OR_PASSWORD")
	}

	err = utils.CompareHashAndPassword(user.Password, password)

	if err != nil {
		return nil, errors.New("INVALID_EMAIL_OR_PASSWORD")
	}

	accessToken, err := utils.GenerateToken(user.ID, false)
	if err != nil {
		return nil, errors.New("ERROR_WHILE_SIGNATURE")
	}
	refreshToken, err := utils.GenerateToken(user.ID, true)
	if err != nil {
		return nil, errors.New("ERROR_WHILE_SIGNATURE")
	}

	res := &response.Token{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	return res, nil
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
