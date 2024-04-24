package ChatService

import (
	"EriChat/middlewares"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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
	middlewares.AuthWebSocket(c, string(jwt), &connection)
	fmt.Println("认证成功！")
	confirmData := utils.ConfirmData()
	persistenceData := utils.PersistenceData()

	self, _ := c.Get("self")
	uid, _ := self.(string)
	if exist, ok := utils.AllConnections.Load(utils.Uid(uid)); ok {
		CloseWebSocket(exist.(*utils.Connection), uid)
	}
	utils.AllConnections.Store(utils.Uid(uid), &connection)
	go connection.ReceiveEvent()
	go connection.SendEvent()
	go func(uid string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			var msg utils.WsMessage
			err = connection.Conn.ReadJSON(&msg)
			if err != nil {
				utils.AllChatRooms.Range(func(key, value any) bool {
					delete(value.(*utils.ChatRoom).Clients, utils.Uid(uid))
					return true
				})
				CloseWebSocket(&connection, uid)
				panic(err)
			}
			switch msg.Type {
			case "message":
				persistenceData <- msg
			case "confirm":
				confirmData <- msg
			case "quitActiveRooms":
				QuitActiveRooms(uid, msg)
			case "quit":
				CloseWebSocket(&connection, uid)
				return
			}
		}
	}(uid)
}

func QuitActiveRooms(uid string, msg utils.WsMessage) {
	for _, qRoom := range msg.QuitRooms {
		if room, ok := utils.AllChatRooms.Load(utils.Cid(qRoom)); ok {
			delete(room.(*utils.ChatRoom).Clients, utils.Uid(uid))
		}
	}
}

func CloseWebSocket(conn *utils.Connection, uid string) {
	conn.ToWS <- utils.WsMessage{
		Type:    "ServerQuit",
		Message: "成功断开WebSocket",
	}
	time.Sleep(time.Second)
	conn.CloseReceive <- true
	conn.CloseSend <- true
	_ = conn.Conn.Close()
	utils.AllConnections.Delete(utils.Uid(uid))
}
