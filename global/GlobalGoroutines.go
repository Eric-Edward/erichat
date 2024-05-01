package global

import (
	"EriChat/models"
	"EriChat/utils"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"strconv"
	"strings"
	"time"
)

type RedisMessagesInfo struct {
	LastUpdate uint
	Latest     uint
}

var timer *time.Timer

var RedisMessages map[string]*RedisMessagesInfo

func InitGlobalGoroutines() {
	timer = time.NewTimer(time.Minute * 5)
	RedisMessages = make(map[string]*RedisMessagesInfo)
	loadInitMessagesIntoRedis()
	messageEventLoop()
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
				redisMessageInfo := RedisMessages[string(msg.Target)]
				redisMessageInfo.Latest++
				redis := utils.GetRedis()
				var message = models.Message{
					ID:       redisMessageInfo.Latest,
					Target:   msg.Target,
					Type:     msg.Type,
					UserName: msg.UserName,
					Uid:      msg.Uid,
					Message:  msg.Message,
				}
				marshal, _ := json.Marshal(message)
				redis.Set(context.Background(), string(msg.Target)+"_"+strconv.Itoa(int(redisMessageInfo.Latest)), marshal, 0)

				msg.ID = redisMessageInfo.Latest
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

	ctx := context.Background()
	redis := utils.GetRedis()
	db := utils.GetMySQLDB()
	err := redis.FlushDB(ctx).Err()
	if err != nil {
		panic(err)
	}

	var cids []string
	db.Model(&models.ChatRoom{}).Select("cid").Find(&cids)
	for _, cid := range cids {
		var messages []models.Message
		db.Table("messages_" + cid).Model(&models.Message{}).Order("id desc").Limit(100).Find(&messages)

		var redisMessageInfo RedisMessagesInfo
		if len(messages) == 0 {
			redisMessageInfo.LastUpdate = 0
			redisMessageInfo.Latest = 0
		} else {
			redisMessageInfo.LastUpdate = messages[0].ID
			redisMessageInfo.Latest = messages[0].ID
		}
		RedisMessages[cid] = &redisMessageInfo

		for _, message := range messages {
			marshal, _ := json.Marshal(message)
			err = redis.Set(ctx, cid+"_"+strconv.Itoa(int(message.ID)), marshal, 0).Err()
			if err != nil {
				panic(err)
			}
		}
	}

}

func persistDataInMysql() {
	ctx := context.Background()
	redis := utils.GetRedis()
	db := utils.GetMySQLDB()

	var cids = maps.Keys(RedisMessages)
	for _, cid := range cids {
		latest := RedisMessages[cid].LastUpdate
		keys, _ := redis.Keys(ctx, cid+"_*").Result()
		var neededInsertMessages []models.Message
		var maxKey uint
		for _, key := range keys {
			split := strings.Split(key, "_")
			if id, _ := strconv.Atoi(split[1]); uint(id) > latest {
				if uint(id) > maxKey {
					maxKey = uint(id)
				}

				var msg models.Message
				message, err := redis.Get(ctx, key).Result()
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal([]byte(message), &msg)
				if err != nil {
					panic(err)
				}
				neededInsertMessages = append(neededInsertMessages, msg)
			} else {
				//redis.Del(ctx, key)
				//这里使用redis.unlink来进行异步释放，减轻redis主线程的压力
				err := redis.Unlink(ctx, key).Err()
				if err != nil {
					panic(err)
				}
			}
		}

		if len(neededInsertMessages) != 0 {
			tx := db.Begin()
			result := tx.Table("messages_" + cid).Create(&neededInsertMessages)
			if result.Error != nil {
				tx.Rollback()
				panic(tx.Error)
			}
			tx.Commit()
			RedisMessages[cid].LastUpdate = maxKey
		}
	}
}
