package ChatService

import (
	"EriChat/global"
	"EriChat/models"
	"EriChat/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"net/http"
	"slices"
	"strconv"
)

func GetMessageByCid(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
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
	var uend uint
	if end != "" {
		e, _ := strconv.Atoi(end)
		uend = uint(e)
	} else {
		uend = global.RedisMessages[cid].LastUpdate + 1
	}
	fmt.Println(uend)

	redis := utils.GetRedis()
	ctx := context.Background()
	var messages []*models.Message
	for i, j := uend, 50; j > 0 && i >= 0; i, j = i-1, j-1 {
		var marshal string
		marshal, err = redis.Get(ctx, cid+"_"+strconv.Itoa(int(i-1))).Result()
		fmt.Println(marshal)
		switch {
		case errors.Is(err, redis2.Nil):
			result, errs := models.GetMessageByCid(cid, i, int(50-uend+i))
			if errs != nil {
				panic(errs)
			}
			slices.Reverse(result)
			messages = append(messages, result...)
			goto next
		case err != nil:
			panic(err)
		default:
			var message models.Message
			err = json.Unmarshal([]byte(marshal), &message)
			if err != nil {
				panic(err)
			}
			messages = append(messages, &message)
		}
	}
next:
	slices.Reverse(messages)
	divider, err := models.GetChatRoomMessageDivider(uid.(string), cid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查询历史信息失败",
			"code":    utils.FailedLoadHistoryMessages,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "获取当前聊天室历史信息成功",
		"code":     utils.Success,
		"messages": messages,
		"divider":  divider,
	})
	return
}
