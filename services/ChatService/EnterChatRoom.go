package ChatService

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// EnterChatRoom 这个函数留作为群聊的进群函数吧
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
	//TODO 在ChatRoomNumber中加上自己
	c.JSON(http.StatusOK, gin.H{
		"message": "成功进入聊天室",
		"code":    utils.Success,
		"cid":     roomInfo.Cid,
		"channel": roomInfo.Channel,
	})
	return
}
