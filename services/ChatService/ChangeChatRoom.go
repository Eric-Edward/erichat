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

	// TODO 后面添加是不是聊天室里的人
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
		chatRoom.Clients = append(chatRoom.Clients, uid.(utils.Uid))
		utils.AllChatRooms[utils.Cid(cid)] = &chatRoom
	} else {
		room.Clients = append(room.Clients, uid.(utils.Uid))
	}

}
