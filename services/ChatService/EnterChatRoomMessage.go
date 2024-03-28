package ChatService

import (
	"EcChat/models"
	"EcChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func CreateChatRoom(c *gin.Context) {
	var chatRoom models.ChatRoom
	err := c.ShouldBind(&chatRoom)
	if err != nil {
		fmt.Println("数据绑定失败！")
		c.JSON(http.StatusOK, gin.H{
			"message": "数据信息绑定失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}
	err = models.CreateChatRoom(chatRoom.Channel)
	if err != nil {
		fmt.Println("创建聊天室失败", err)
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

func EnterChatRoom(c *gin.Context) {
	var chatRoom models.ChatRoom
	err := c.ShouldBind(&chatRoom)
	if err != nil {
		fmt.Println("数据绑定失败！")
		c.JSON(http.StatusOK, gin.H{
			"message": "数据信息绑定失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}
	roomInfo, err := models.GetChatRoomByCid(chatRoom.Cid)
	if err != nil {
		fmt.Println("聊天室不存在")
		c.JSON(http.StatusOK, gin.H{
			"message": "当前查询的聊天室不存在",
			"code":    utils.FailedFindChatRoom,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功进入聊天室",
		"code":    utils.Success,
		"cid":     roomInfo.Cid,
		"channel": roomInfo.Channel,
	})
	return
}

func SendMessage(c *gin.Context) {

}
