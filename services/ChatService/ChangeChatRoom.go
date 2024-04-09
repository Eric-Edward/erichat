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

	room, ok := utils.AllChatRooms[utils.Cid(cid)]
	if !ok {
		chatRoom := utils.ChatRoom{
			Cid:     cid,
			Clients: make([]utils.Uid, 0),
		}
		chatRoom.Clients = append(chatRoom.Clients, utils.Uid(uid.(string)))
		utils.AllChatRooms[utils.Cid(cid)] = &chatRoom
	} else {
		room.Clients = append(room.Clients, utils.Uid(uid.(string)))
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "切换当前聊天室成功",
		"code":    utils.Success,
	})
	return
}
