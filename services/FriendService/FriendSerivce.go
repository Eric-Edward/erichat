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
	if username == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "查找的用户名不应为空",
			"code":    utils.FailedFindUser,
		})
		return
	}
	clients, err := models.GetAllClientsByUserName(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查找用户失败",
			"code":    utils.FailedFindClients,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查找成功",
		"code":    utils.Success,
		"clients": clients,
	})
	return
}

func AddFriend(c *gin.Context) {
	var relationShip models.RelationShipApply
	err := c.ShouldBind(&relationShip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "数据绑定失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}

	fmt.Println(relationShip.Apply, relationShip.Applied)

	result, err := models.ApplyRelationShip(relationShip.Apply, relationShip.Applied, relationShip.Group)
	if err != nil || !result {
		c.JSON(http.StatusOK, gin.H{
			"message": "添加朋友失败",
			"code":    utils.FailedAddFriends,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "申请添加朋友成功",
		"code":    utils.Success,
	})
	return
}

func GetAllApplyByUid(c *gin.Context) {
	uid := c.Query("uid")
	applies, err := models.GetRelationShipApplyByUid(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "获取朋友申请失败",
			"code":    utils.FailedLoadApplies,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取朋友申请成功",
		"code":    utils.Success,
		"applies": applies,
	})
	return
}

func GetGroupByUid(c *gin.Context) {
	uid := c.Query("uid")
	groups, err := models.GetAllGroupByUid(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "获取用户分组失败",
			"code":    utils.FailedLoadGroups,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "加载分组成功",
		"code":    utils.Success,
		"groups":  groups,
	})
	return
}

func AddGroup(c *gin.Context) {
	var group models.Group
	err := c.ShouldBind(&group)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "数据绑定失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}
	result, err := models.AddGroup(group.Uid, group.GroupName)
	if err != nil || !result {
		c.JSON(http.StatusOK, gin.H{
			"message": "添加分组失败！",
			"code":    utils.FailedAddGroup,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "添加分组成功",
		"code":    utils.Success,
	})
}

func AddRelationShip(c *gin.Context) {
	var relationShipApply models.RelationShipApply
	err := c.ShouldBind(&relationShipApply)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "数据绑定失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}

	result, err := models.HandleRelationShipApply(relationShipApply.Apply, relationShipApply.Applied, relationShipApply.Group)
	if err != nil || !result {
		c.JSON(http.StatusOK, gin.H{
			"message": "添加朋友关系失败",
			"code":    utils.FailedAddFriends,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "添加朋友关系成功",
		"code":    utils.Success,
	})
	return

}

func DeleteRelationShipApply(c *gin.Context) {
	apply := c.Query("Apply")
	applied := c.Query("Applied")
	result, err := models.DropRelationShipApply(apply, applied)
	if err != nil || !result {
		c.JSON(http.StatusOK, gin.H{
			"message": "删除用户请求失败",
			"code":    utils.FailedDropApply,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "拒绝用户申请成功",
		"code":    utils.Success,
	})
	return
}
