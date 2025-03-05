package main

import (
	"fmt"
	"net/http"

	"authz/internal/controllers"
	"authz/internal/initializers"
	"authz/internal/login"
	"authz/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func mux_http_router() {
    router := mux.NewRouter()
    router.HandleFunc("/login", login.LoginHandler).Methods("POST")
    router.HandleFunc("/protected", login.ProtectedHandler).Methods("GET")

    fmt.Println("Server is running on port 8000")
    err := http.ListenAndServe(":8000", router)    
    if err != nil {
        fmt.Println("Could not start the server", err)
    } 
    fmt.Println("Server is running on port 8000")
}


func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func gin_http_router() {
    router := gin.Default()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    
    api := router.Group("/api")
    {
        //endpoint : /api/auth/signup and /api/auth/login
        auth_api := api.Group("/auth")
        {
            auth_api.POST("/signup", controllers.CreateUser)
            auth_api.POST("/login", controllers.Login)
        }

        //endpoint : /api/user/profile and /api/user/validate
        user_api := router.Group("/user")
        {
           user_api.GET("/profile", middleware.CheckAuthentication, controllers.GetUserProfile)
           user_api.GET("/validate", controllers.ValidateToken)
        }
    }
    
    router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{"Message":"Invalid API endpoint"}) })
    router.Run(":8000")

}


func main() {
    gin_http_router()
}