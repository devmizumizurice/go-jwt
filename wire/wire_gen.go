// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/devmizumizurice/go-jwt/controllers"
	"github.com/devmizumizurice/go-jwt/repositories"
	"github.com/devmizumizurice/go-jwt/services"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitializeAuthController(db *gorm.DB) controllers.AuthController {
	userRepositoryInterface := repositories.NewUserRepository(db)
	userServiceInterface := services.NewUserService(userRepositoryInterface)
	authController := controllers.NewAuthController(userServiceInterface)
	return authController
}