package ChatService

import (
	"EriChat/middlewares"
	"EriChat/models"
	"EriChat/utils"
	"context"
	"encoding/json"
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

	self, _ := c.Get("self")
	uid, _ := self.(string)
	utils.AllConnections.Store(utils.Uid(uid), &connection)
	go connection.ReceiveEvent()
	go connection.SendEvent()
	go func() {
		for {
			var msg utils.WsMessage
			err = connection.Conn.ReadJSON(&msg)
			if err != nil {
				_ = connection.Conn.Close()
			}
			if msg.Type == "message" {
				connection.FromWS <- msg
				HandleMessage(msg)
			} else if msg.Type == "quitActiveRooms" {
				for _, qRoom := range msg.QuitRooms {
					if room, ok := utils.AllChatRooms.Load(utils.Cid(qRoom)); ok {
						delete(room.(*utils.ChatRoom).Clients, utils.Uid(uid))
						if len(room.(*utils.ChatRoom).Clients) == 0 {
							models.AddMessageToTable(room.(*utils.ChatRoom))
						}
					}
				}
			} else {
				connection.CloseReceive <- true
				connection.CloseSend <- true
				utils.AllConnections.Delete(utils.Uid(uid))
				_ = connection.Conn.Close()
				return
			}
		}
	}()
}

func HandleMessage(msg utils.WsMessage) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	value, _ := utils.AllChatRooms.Load(msg.Target)
	room := value.(*utils.ChatRoom)
	if len(room.Clients) != 1024 {
		room.Message = append(room.Message, msg)
	} else {
		redis := utils.GetRedis()
		var ctx = context.Background()
		var messages []utils.WsMessage
		err := json.Unmarshal([]byte(redis.Get(ctx, string(msg.Target)).Val()), &messages)
		if err != nil {
			panic(err)
		}

		messages = append(messages, room.Message...)
		marshal, err := json.Marshal(messages)
		_, err = redis.Set(ctx, string(msg.Target), marshal, time.Hour).Result()
		if err != nil {
			panic(err)
		}

		//将信息放到redis中后，将内存释放
		room.Message = room.Message[:0]
	}

}
