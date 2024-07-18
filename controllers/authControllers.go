package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/devmizumizurice/go-jwt/models"
	"github.com/devmizumizurice/go-jwt/services"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	RefreshToken(c *gin.Context)
	Validate(c *gin.Context)
}

type authController struct {
	userService services.UserServiceInterface
}

func NewAuthController(userService services.UserServiceInterface) AuthController {
	return &authController{userService: userService}
}

func (ac *authController) SignUp(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INCORRECT_BODY"})
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := ac.userService.CreateUser(user)
	if err != nil {
		if err.Error() == "EMAIL_ALREADY_EXISTS" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (ac *authController) SignIn(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INCORRECT_BODY"})
		return
	}

	token, err := ac.userService.VerifyUserEmail(req.Email, req.Password)

	if err != nil {
		if err.Error() == "INVALID_EMAIL_OR_PASSWORD" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	accessTokenLifeMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_VALIDATE_MINUTES"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}

	refreshTokenLifeDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_VALIDATE_DAYS"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}

	c.SetCookie("access_token", token.AccessToken, accessTokenLifeMinutes*60, "", "", false, true)
	c.SetCookie("refresh_token", token.RefreshToken, refreshTokenLifeDays*24*60*60, "", "", false, true)
	c.JSON(http.StatusNoContent, gin.H{})
}

func (ac *authController) RefreshToken(c *gin.Context) {
	// TODO:
}

func (ac *authController) Validate(c *gin.Context) {
	sub, _ := c.Get("sub")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("I'm logged in %s", sub),
	})
}
