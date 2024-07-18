package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/devmizumizurice/go-jwt/models"
	"github.com/devmizumizurice/go-jwt/services"
	"github.com/gin-gonic/gin"
)

type AuthControllerInterface interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authController struct {
	userService services.AuthServiceInterface
}

func NewAuthController(userService services.AuthServiceInterface) AuthControllerInterface {
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
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "MISSING_TOKEN"})
		return
	}
	newToken, err := ac.userService.RefreshToken(token)

	if err != nil {
		if err.Error() == "USER_NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "EXPIRED_TOKEN" || err.Error() == "UNAUTHORIZED" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	accessTokenLifeMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_VALIDATE_MINUTES"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshTokenLifeDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_VALIDATE_DAYS"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.SetCookie("access_token", newToken.AccessToken, accessTokenLifeMinutes*60, "", "", false, true)
	c.SetCookie("refresh_token", newToken.RefreshToken, refreshTokenLifeDays*24*60*60, "", "", false, true)
	c.JSON(http.StatusNoContent, gin.H{})

}
