package main

import (
	"go-restfull-api/config"
	"go-restfull-api/handler"
	"go-restfull-api/repository"
	"go-restfull-api/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitDB()

	router := gin.Default()

	userRepo := repository.NewUserRepositoryPostgres(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(*userService)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	{
		userRouter := router.Group("/user")
		userRouter.POST("/auth/signup", userHandler.HandleSignup)
		userRouter.POST("/auth/signin", userHandler.HandleSignin)
		userRouter.GET("/auth/me", userHandler.HandleGetProfile)
	}

	router.Run()
}
