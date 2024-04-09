package ChatService

import (
	"EriChat/models"
	"EriChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateChatRoom(c *gin.Context) {
	uid, _ := c.Get("self")
	chatRoom := make(map[string]interface{})
	err := c.BindJSON(&chatRoom)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "获取聊天室成员失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}
	clients := chatRoom["clients"].([]interface{})
	clients = append(clients, uid.(string))
	chatRoomName := chatRoom["chatRoomName"].(string)
	result, err := models.CreateChatRoom(chatRoomName, clients)
	if err != nil || !result {
		c.JSON(http.StatusOK, gin.H{
			"message": "创建聊天室失败",
			"code":    utils.FailedCreateChatRoom,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "创建聊天室成功",
		"code":    utils.Success,
	})
	return
}
