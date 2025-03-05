package middleware

import (
	"fmt"
	"authz/internal/initializers"
	"authz/internal/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyToken(c *gin.Context) (models.User, error) {
    authHeader := c.GetHeader("Authorization")

    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
        return models.User{}, fmt.Errorf("authorization header is required")
    }

    authTokens := strings.Split(authHeader, " ")
    if len(authTokens) != 2 || authTokens[0] != "Bearer" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is invalid"})
        return models.User{}, fmt.Errorf("authorization header is invalid")
    }

    tokenString := authTokens[1]

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }   
        return []byte(os.Getenv("SECRET")), nil
    })

    if err != nil || !token.Valid {
        fmt.Printf("Invalid Token ERROR: [%s]\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        return models.User{}, fmt.Errorf("Invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        return models.User{}, fmt.Errorf("Invalid token")
    }

    if float64(time.Now().Unix()) > claims["exp"].(float64) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
        return models.User{}, fmt.Errorf("Token expired")
    }

    var user models.User
    initializers.DB.Where("id=?", claims["id"]).First(&user)
    if user.ID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
        return models.User{}, fmt.Errorf("User not found")
    }

    return user, nil
}

func CheckAuthentication(c *gin.Context) {
    user, err := VerifyToken(c)
    if err != nil {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }
    c.Set("currentUser", user)
    c.Next()
}

func CheckAuth(c *gin.Context){
    authHeader := c.GetHeader("Authorization")

    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    authTokens := strings.Split(authHeader, " ")
    if len(authTokens) != 2 || authTokens[0] != "Bearer" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is invalid"})
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    tokenString := authTokens[1]

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }   
        return []byte(os.Getenv("SECRET")), nil
    })

    if err != nil || !token.Valid {
        fmt.Printf("Invalid Token ERROR: [%s]\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    if float64(time.Now().Unix()) > claims["exp"].(float64) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    var user models.User
    initializers.DB.Where("id=?", claims["id"]).First(&user)
    if user.ID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    c.Set("currentUser", user)
    c.Next()
}