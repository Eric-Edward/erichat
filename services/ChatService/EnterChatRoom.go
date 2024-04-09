package ChatService

import (
	"EriChat/middlewares"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CreateWebSocketConn EnterChatRoom 这个函数留作为群聊的进群函数
func CreateWebSocketConn(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "创建websocket失败",
			"code":    utils.FailedCreateWebSocket,
		})
		return
	}
	connection := utils.Connection{
		Conn:         ws,
		FromWS:       make(chan utils.WsMessage),
		ToWS:         make(chan utils.WsMessage),
		CloseReceive: make(chan bool),
		CloseSend:    make(chan bool),
	}

	_, jwt, err := connection.Conn.ReadMessage()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "读取WebSocket失败",
			"code":    utils.FailedCreateWebSocket,
		})
		_ = connection.Conn.Close()
		return
	}
	middlewares.AuthWebSocket(c, string(jwt))

	self, _ := c.Get("self")
	uid, _ := self.(string)
	utils.AllConnections[utils.Uid(uid)] = &connection
	go connection.ReceiveEvent()
	go connection.SendEvent()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("从客户端接受消息失败:", err)
			}
		}()
		for {
			var msg utils.WsMessage
			err = connection.Conn.ReadJSON(&msg)
			if err != nil {
				_ = connection.Conn.Close()
			}
			if msg.Type == "message" {
				connection.FromWS <- msg
			} else if msg.Type == "quitActiveRooms" {
				for _, room := range msg.QuitRooms {
					delete(utils.AllChatRooms, utils.Cid(room))
				}
			} else {
				connection.CloseReceive <- true
				connection.CloseSend <- true
				_ = connection.Conn.Close()
				delete(utils.AllConnections, utils.Uid(uid))
				return
			}
		}
	}()
}
