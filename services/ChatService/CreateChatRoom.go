package ChatService

import "github.com/gin-gonic/gin"

func CreatePeerChatRoom(c *gin.Context) {
	//	u1 := c.Param("u1")
	//	u2 := c.Param("u2")
	//	channel := c.Param("channel")
	//	var cid string
	//	cid, err := models.CreatePeerChatRoom(channel, u1, u2)
	//	if err != nil {
	//		fmt.Println("创建聊天室失败", err)
	//		c.JSON(http.StatusOK, gin.H{
	//			"message": "创建聊天室失败",
	//			"code":    utils.FailedCreateChatRoom,
	//		})
	//		return
	//	}
	//
	//	//conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	//
	//	//这里我们先记录这个被动用户的信息，当用户上线后，会先进行chat的搜索，如果搜索到没有加入的就被动的创建一个socket连接
	//	//socket.NewWSClient(u1, cid, channel, conn)
	//
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "创建聊天室成功",
	//		"code":    utils.Success,
	//	})
	return
}
