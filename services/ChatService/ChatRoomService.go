package ChatService

import (
	"EriChat/models"
	"EriChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllChatGroup(c *gin.Context) {
	uid, _ := c.Get("self")
	rooms := models.GetAllChatGroupByUid(uid.(string))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取当前用户的全部聊天时",
		"code":    utils.Success,
		"rooms":   rooms,
	})
	return
}

func GetChatRoomInfoByCid(c *gin.Context) {
	cid := c.Query("cid")
	room, err := models.GetChatRoomByCid(cid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询当前聊天室信息失败",
			"code":    utils.FailedFindChatRoom,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询聊天室信息成功",
		"code":    utils.Success,
		"room":    room,
	})
	return
}
