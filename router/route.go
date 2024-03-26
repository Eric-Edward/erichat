package router

import (
	"EcChat/docs"
	"EcChat/middlewares"
	"EcChat/services/UserService"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func Route() *gin.Engine {
	ginServer := gin.Default()

	ginServer.Use(middlewares.Cors())

	docs.SwaggerInfo.BasePath = ""
	ginServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routeLogin := ginServer.Group("/login")
	{
		routeLogin.POST("", UserService.LoginHandler)
	}

	routeRegister := ginServer.Group("/register")
	{
		routeRegister.GET("", UserService.UsernameIsRegistered)
		routeRegister.POST("", UserService.Register)
	}

	routeComplete := ginServer.Group("/complete")
	{
		routeComplete.GET("", UserService.GetUserInfo)
		routeComplete.POST("", UserService.CompleteUserInfo)
	}

	return ginServer
}
