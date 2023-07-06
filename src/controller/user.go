package controller

import (
	"net/http"
	"omni-backend-service/src/model"
	"omni-backend-service/src/service"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	COOKIE_DURATION = int(1 * time.Hour)
)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) AuthController {
	return &userController{userService: userService}
}

func (u *userController) Register(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusCreated, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	user.ID = c.GetString("userID")
	err := u.userService.CreateUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	jwt := c.GetString("jwt")
	c.SetCookie("jwt", jwt, COOKIE_DURATION, "/", "", false, false)

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (u *userController) Login(c *gin.Context) {
	userID := c.GetString("userID")
	_, err := u.userService.GetUserByID(userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	jwt := c.GetString("jwt")
	c.SetCookie("jwt", jwt, COOKIE_DURATION, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (u *userController) GetUser(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := u.userService.GetUserByID(string(userID))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "status": "success"})
}

func (u *userController) UpdateUser(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusCreated, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	user.ID = c.GetString("userID")
	err := u.userService.UpdateUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}
