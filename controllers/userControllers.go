package controllers

import (
	"fmt"
	"net/http"

	"github.com/devmizumizurice/go-jwt/services"
	"github.com/gin-gonic/gin"
)

type UserControllerInterface interface {
	Validate(c *gin.Context)
}

type userController struct {
	userService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) UserControllerInterface {
	return &userController{userService: userService}
}

func (ac *userController) Validate(c *gin.Context) {
	sub, _ := c.Get("sub")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("I'm logged in %s", sub),
	})
}
