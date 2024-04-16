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
