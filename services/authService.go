package services

import (
	"errors"
	"time"

	"github.com/devmizumizurice/go-jwt/models"
	"github.com/devmizumizurice/go-jwt/models/response"
	"github.com/devmizumizurice/go-jwt/repositories"
	"github.com/devmizumizurice/go-jwt/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceInterface interface {
	CreateUser(user *models.User) (*response.User, error)
	VerifyUserEmail(email string, password string) (*response.Token, error)
	RefreshToken(refreshToken string) (*response.Token, error)
}

type authService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewAuthService(userRepository repositories.UserRepositoryInterface) AuthServiceInterface {
	return &authService{userRepository: userRepository}
}

func (s *authService) CreateUser(user *models.User) (*response.User, error) {
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

func issueToken(userId string) (*response.Token, error) {
	accessToken, err := utils.GenerateToken(userId, false)
	if err != nil {
		return nil, errors.New("ERROR_WHILE_SIGNATURE")
	}
	refreshToken, err := utils.GenerateToken(userId, true)
	if err != nil {
		return nil, errors.New("ERROR_WHILE_SIGNATURE")
	}

	return &response.Token{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (s *authService) VerifyUserEmail(email string, password string) (*response.Token, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("INVALID_EMAIL_OR_PASSWORD")
	}

	err = utils.CompareHashAndPassword(user.Password, password)

	if err != nil {
		return nil, errors.New("INVALID_EMAIL_OR_PASSWORD")
	}

	return issueToken(user.ID)
}

func (s *authService) RefreshToken(refreshToken string) (*response.Token, error) {
	parsedToken, err := utils.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.New("UNAUTHORIZED")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if ok && parsedToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("EXPIRED_TOKEN")
		}
	}
	sub := claims["sub"].(string)
	user, err := s.userRepository.FindByID(sub)
	if err != nil {
		return nil, errors.New("USER_NOT_FOUND")
	}

	return issueToken(user.ID)
}
