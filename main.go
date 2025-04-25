package main

import (
	"go-restfull-api/config"
	"go-restfull-api/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitDB()

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	{
		userRouter := router.Group("/user")
		userRouter.POST("/auth/signup", handler.HandleSignup)
		userRouter.POST("/auth/signin", handler.HandleSignin)
		userRouter.GET("/auth/me", handler.HandleGetProfile)
	}

	router.Run()
}
