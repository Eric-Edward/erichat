package ChatService

import (
	"EriChat/models"
	"EriChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChangeChatRoom(c *gin.Context) {
	cid := c.Query("cid")
	uid, _ := c.Get("self")

	isMember, err := models.IsChatRoomMember(cid, uid.(string))
	if err != nil || !isMember {
		c.JSON(http.StatusOK, gin.H{
			"message": "非本聊天室成员",
			"code":    utils.NotChatRoomMember,
		})
		return
	}

	room, ok := utils.AllChatRooms.Load(utils.Cid(cid))
	if !ok {
		chatRoom := utils.ChatRoom{
			Cid:     cid,
			Clients: make(map[utils.Uid]utils.Uid),
		}
		chatRoom.Clients[utils.Uid(uid.(string))] = utils.Uid(uid.(string))
		utils.AllChatRooms.Store(utils.Cid(cid), &chatRoom)
	} else {
		room.(*utils.ChatRoom).Clients[utils.Uid(uid.(string))] = utils.Uid(uid.(string))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "切换当前聊天室成功",
		"code":    utils.Success,
	})
	return
}
