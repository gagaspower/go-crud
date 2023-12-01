package main

import (
	"learn-go/config"
	authcontroller "learn-go/controller/authController"
	"learn-go/controller/productController"
	"learn-go/controller/userController"
	"learn-go/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.ConnectDatabase()

	v1 := r.Group("/api/v1")
	{
		v1.POST("login", authcontroller.Login)

		protected := r.Group("/api/v1/admin")
		protected.Use(middlewares.JwtAuthMiddleware())

		protected.GET("products", productController.Index)
		protected.GET("product/:Id", productController.Show)
		protected.POST("product", productController.Create)
		protected.PUT("product/:Id", productController.Update)
		protected.DELETE("product/:Id", productController.Destroy)

		// route for user controller
		protected.GET("users", userController.Index)
		protected.POST("user", userController.Create)
		protected.GET("user/:id", userController.Show)
		protected.PUT("user/:id", userController.Update)
		protected.DELETE("user/:id", userController.Destroy)
	}

	r.Run()
}
