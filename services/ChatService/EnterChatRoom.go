package ChatService

import (
	"EriChat/utils"
	"encoding/json"
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
	uid, _ := c.Get("self")
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "创建websocket失败",
			"code":    utils.FailedCreateWebSocket,
		})
		return
	}
	connection := utils.Connection{
		Conn:   ws,
		FromWS: make(chan utils.WsMessage),
		ToWS:   make(chan utils.WsMessage),
	}
	utils.AllConnections[uid.(utils.Uid)] = &connection
	go connection.EventLoop()
	go func() {
		for {
			_, p, err2 := connection.Conn.ReadMessage()
			if err2 != nil {
				_ = connection.Conn.Close()
			}
			var msg utils.WsMessage
			err = json.Unmarshal(p, &msg)
			if err != nil {
				continue
			}
			connection.FromWS <- msg
		}
	}()
}
