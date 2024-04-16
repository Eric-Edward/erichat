package ChatService

import (
	"EriChat/models"
	"EriChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetMessageByCid(c *gin.Context) {
	cid := c.Query("cid")
	end := c.Query("end")
	uid, _ := c.Get("self")
	member, err := models.IsChatRoomMember(cid, uid.(string))
	if err != nil || !member {
		c.JSON(http.StatusOK, gin.H{
			"message": "非当前用户成员",
			"code":    utils.FailedFindUser,
		})
		return
	}

	uend, _ := strconv.Atoi(end)
	messages, err := models.GetMessageByCid(cid, uint(uend))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "获取聊天室历史信息失败",
			"code":    utils.FailedLoadHistoryMessages,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "获取当前聊天室历史信息成功",
		"code":     utils.Success,
		"messages": messages,
	})
	return
}
