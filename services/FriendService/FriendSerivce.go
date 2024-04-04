package FriendService

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllFriends(c *gin.Context) {
	uid := c.Query("uid")
	relations, err := models.GetRelationShipByUid(uid)
	if err != nil {
		fmt.Println("获取当前用户朋友列表失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "获取当前用户朋友列表失败",
			"code":    utils.FailedLoadFriends,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取当前用户列表成功",
		"code":    utils.Success,
		"friends": relations,
	})
	return
}

func GetAllClientByUserName(c *gin.Context) {
	username := c.Query("username")
	db := utils.GetMySQLDB()
}
