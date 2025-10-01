package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shravanskie/vegetable-backend/internal/db"
	"github.com/shravanskie/vegetable-backend/internal/handlers"
	"github.com/shravanskie/vegetable-backend/internal/middleware"
	"github.com/shravanskie/vegetable-backend/internal/models"
	"github.com/shravanskie/vegetable-backend/internal/services"
)

func main() {
	// Connect to MySQL
	db.Connect()
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Vegetable{})
	// Initialize UserService
	userService := services.NewUserService("YOUR_GOOGLE_CLIENT_ID")

	// Initialize Handlers
	signupHandler := handlers.NewSignupHandler(userService)
	vegtableHandler := handlers.NewVegetableHandler(services.NewVegetableService())
	r := gin.Default()
	pub := r.Group("/api")
	{
		pub.POST("/signup", signupHandler.Signup)
		pub.POST("/google-signup", signupHandler.GoogleSignup)
		pub.POST("/login", signupHandler.Login)
		pub.GET("/validate-token", signupHandler.ValidateToken)
	}
	r.Static("/images", "./uploads/vegetables")
	prot := r.Group("/api")
	prot.Use(middleware.AuthMiddleware())
	{
		prot.POST("/vegetables", vegtableHandler.AddVegetable)  // add vegetable with image
		prot.GET("/vegetables", vegtableHandler.ListVegetables) // list vegetables
	}
	//api.Use(middleware.AuthMiddleware())
	r.Run(":8080")
}
