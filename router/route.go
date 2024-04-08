package router

import (
	"EriChat/middlewares"
	"EriChat/services/ChatService"
	"EriChat/services/FriendService"
	"EriChat/services/UserService"
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

	routeSocket := ginServer.Group("/enter")
	{
		routeSocket.GET("", ChatService.CreateWebSocketConn)
	}

	routeComplete := ginServer.Group("/complete", middlewares.Authorization())
	{
		routeComplete.GET("", UserService.GetUserInfo)
		routeComplete.POST("", UserService.CompleteUserInfo)
	}

	routeMessage := ginServer.Group("/chat", middlewares.Authorization())
	{
		routeMessage.POST("/createPeer", ChatService.CreatePeerChatRoom)
		routeMessage.POST("/message", ChatService.ReceiveMessage)
		routeMessage.GET("/chatRoom", ChatService.GetAllChatRoom)
		routeMessage.GET("/changeChatRoom", ChatService.ChangeChatRoom)
	}

	routeFriends := ginServer.Group("/friends", middlewares.Authorization())
	{
		routeFriends.GET("", FriendService.GetAllFriends)
		routeFriends.POST("", FriendService.AddFriend)
		routeFriends.GET("/group", FriendService.GetGroupByUid)
		routeFriends.POST("/group", FriendService.AddGroup)
		routeFriends.GET("/apply", FriendService.GetAllApplyByUid)
		routeFriends.POST("/apply", FriendService.AddRelationShip)
		routeFriends.DELETE("/apply", FriendService.DeleteRelationShipApply)
	}

	routeClient := ginServer.Group("/clients", middlewares.Authorization())
	{
		routeClient.GET("", FriendService.GetAllClientByUserName)
	}

	return ginServer
}
