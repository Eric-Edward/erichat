package router

import (
	"EcChat/middlewares"
	"EcChat/services/ChatService"
	"EcChat/services/UserService"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	ginServer := gin.Default()

	ginServer.Use(middlewares.Cors())

	routeLogin := ginServer.Group("/login")
	{
		routeLogin.POST("", UserService.LoginHandler)
	}

	routeRegister := ginServer.Group("/register")
	{
		routeRegister.GET("", UserService.UsernameIsRegistered)
		routeRegister.POST("", UserService.Register)
	}

	routeComplete := ginServer.Group("/complete", middlewares.Authorization())
	{
		routeComplete.GET("", UserService.GetUserInfo)
		routeComplete.POST("", UserService.CompleteUserInfo)
	}

	routeMessage := ginServer.Group("/chat")
	{
		routeMessage.GET("/enter", ChatService.EnterChatRoom)
	}

	return ginServer
}
