package main

import (
	"user-management-api/config"
	"user-management-api/internal/handler"
	"user-management-api/internal/repository"
	"user-management-api/internal/service"
	"user-management-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()
	r := gin.Default()

	repo := repository.NewUserRepository(db)
	usecase := service.NewUserUsecase(repo)
	handler := handler.NewUserHandler(usecase)

	api := r.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.GET("/profile", middleware.JWTAuthMiddleware(), handler.Profile)
		api.GET("/users", middleware.JWTAuthMiddleware(), handler.Users)
	}

	r.Run()
}
