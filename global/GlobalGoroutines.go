package global

import (
	"EriChat/models"
	"EriChat/utils"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"
)

var timer *time.Timer

type RedisMessage struct {
	Cursor   uint
	Messages []models.Message
}

func InitGlobalGoroutines() {
	loadInitMessagesIntoRedis()
	messageEventLoop()
	timer = time.NewTimer(time.Minute)
}

func messageEventLoop() {
	persistenceData := utils.PersistenceData()
	confirmData := utils.ConfirmData()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			select {
			case msg := <-persistenceData:
				var redisMessage RedisMessage
				redis := utils.GetRedis()
				result, err := redis.Get(context.Background(), string(msg.Target)).Result()
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal([]byte(result), &redisMessage)
				if err != nil {
					panic(err)
				}
				redisMessage.Cursor++
				var message = models.Message{
					Target:   msg.Target,
					Type:     msg.Type,
					Message:  msg.Message,
					UserName: msg.UserName,
					Uid:      msg.Uid,
				}
				message.ID = redisMessage.Cursor
				redisMessage.Messages = append(redisMessage.Messages, message)
				marshal, err := json.Marshal(redisMessage)
				if err != nil {
					panic(err)
				}
				err = redis.Set(context.Background(), string(msg.Target), marshal, 0).Err()
				if err != nil {
					panic(err)
				}

				msg.ID = redisMessage.Cursor
				connection, _ := utils.AllConnections.Load(msg.Uid)
				connection.(*utils.Connection).FromWS <- msg
			case msg := <-confirmData:
				db := utils.GetMySQLDB()
				tx := db.Begin()
				result := tx.Model(&models.ChatRoomMember{}).Where("cid=? and uid=? and record<?", msg.Target, msg.Uid, msg.Message).Update("record", msg.Message)
				if result.Error != nil {
					tx.Rollback()
					panic(tx.Error)
				}
				tx.Commit()
			case <-timer.C:
				persistDataInMysql()
			}
		}
	}()
}

func loadInitMessagesIntoRedis() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	redis := utils.GetRedis()
	db := utils.GetMySQLDB()

	var cids []string
	db.Model(&models.ChatRoom{}).Select("cid").Find(&cids)
	for _, cid := range cids {
		var messages []models.Message
		db.Table("messages_" + cid).Model(&models.Message{}).Order("id desc").Limit(100).Find(&messages)
		slices.Reverse(messages)
		var redisMessage RedisMessage
		if len(messages) == 0 {
			redisMessage = RedisMessage{
				Cursor:   0,
				Messages: messages,
			}
		} else {
			redisMessage = RedisMessage{
				Cursor:   messages[len(messages)-1].ID,
				Messages: messages,
			}
		}
		rmsg, _ := json.Marshal(&redisMessage)
		err := redis.Set(context.Background(), cid, rmsg, 0).Err()
		if err != nil {
			fmt.Println("这里是添加信息到redis")
			panic(err)
		}
	}
}

func persistDataInMysql() {
	redis := utils.GetRedis()
	cids, err := redis.Keys(context.Background(), "*").Result()
	if err != nil {
		fmt.Println("从redis中获取cid失败")
		panic(err)
	}
	for _, cid := range cids {
		MockTest("Eric", cid)
	}
}
