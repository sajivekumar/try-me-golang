package controllers

import (
	"authz/internal/initializers"
	"authz/internal/middleware"
	"authz/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


func CreateUser(c *gin.Context) {

	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: authInput.Username,
		Password: string(passwordHash),
	}

	initializers.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func Login(c *gin.Context) {

	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

    /*simpleClaims := jwt.MapClaims{
        "id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}*/

    userClaims := jwt.MapClaims{
        "id":  userFound.ID,
        "name":  userFound.Username,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	// generateSimpleToken := jwt.NewWithClaims(jwt.SigningMethodHS256, simpleClaims)
    generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func ValidateToken(c *gin.Context) {

    user, err := middleware.VerifyToken(c)
    if err != nil {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }
    
	if (user == models.User{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

    c.JSON(200, gin.H{
		"status": "Token Valid",
	})
}


func GetUserProfile(c *gin.Context) {
	user, _ := c.Get("currentUser")

	c.JSON(200, gin.H{
		"user": user,
	})
}