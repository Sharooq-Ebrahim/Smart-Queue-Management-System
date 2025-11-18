package controller

import (
	"log"
	"net/http"
	"smart-queue/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func (auth *AuthController) Register(c *gin.Context) {

	var userInfo model.User

	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Input"})
		return
	}

	var existing model.User

	err := auth.DB.Where("user_name = ?", userInfo.UserName).First(&existing).Error

	log.Println("error", err)
	if err == nil {
		// record found
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)

	user := model.User{
		UserName: userInfo.UserName,
		Password: string(hashed),
	}

	if err := auth.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unble to create user"})
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "User Create sucessfully"})

}

func (auth *AuthController) Login(c *gin.Context) {

	type AuthRequest struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	var LoginRequest AuthRequest

	if err := c.ShouldBindJSON(&LoginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var userInfo model.User

	if err := auth.DB.Where("user_name = ?", LoginRequest.UserName).First(&userInfo).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(LoginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password missmatches"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "User login successfully"})

}
