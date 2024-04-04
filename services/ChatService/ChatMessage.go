package ChatService

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ReceiveMessage(c *gin.Context) {
	var chatRoom models.ChatRoom
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	defer func() { _ = conn.Close() }()
	if err != nil {
		fmt.Println("生成websocket连接失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "生成socket连接失败",
			"code":    utils.FailedGenerateSocket,
		})
		return
	}
	err = c.ShouldBind(&chatRoom)
	if err != nil {
		fmt.Println("数据绑定失败！")
		c.JSON(http.StatusOK, gin.H{
			"message": "数据信息绑定失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}
	_, p, err := conn.ReadMessage()
	if err != nil {

		fmt.Println("ws读取数据失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "从websocket读取数据失败",
			"code":    utils.FailedReadMessage,
		})
		return
	}
	utils.Publish(c, chatRoom.Channel, string(p))
}

func SendMessage(c *gin.Context) {
	channel := c.Param("channel")
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	defer func() { _ = conn.Close() }()
	if err != nil {
		fmt.Println("生成websocket连接失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "生成socket连接失败",
			"code":    utils.FailedGenerateSocket,
		})
		return
	}
	utils.Subscribe(c, channel)
}

func GetAllChatRoom(c *gin.Context) {
	uid := c.Query("uid")
	rooms := models.GetAllChatRoomByUid(uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "获取当前用户的全部聊天时",
		"code":    utils.Success,
		"rooms":   rooms,
	})
	return
}
