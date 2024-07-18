package main

import (
	"github.com/devmizumizurice/go-jwt/initializers"
	"github.com/devmizumizurice/go-jwt/middleware"
	"github.com/devmizumizurice/go-jwt/wire"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.SetUpDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()

	db := initializers.GetDB()
	authController := wire.InitializeAuthController(db)
	userController := wire.InitializeUserController(db)

	r.POST("/signup", authController.SignUp)
	r.POST("/signin", authController.SignIn)
	r.POST("/refresh", authController.RefreshToken)
	r.GET("/validate", middleware.RequireToken, userController.Validate)

	r.Run()

}
