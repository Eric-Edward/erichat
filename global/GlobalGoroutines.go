package global

import (
	"EriChat/models"
	"EriChat/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var timer *time.Timer

var redisMessages map[string]uint

func InitGlobalGoroutines() {
	timer = time.NewTimer(time.Minute * 2)
	redisMessages = make(map[string]uint)
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
				id := redisMessages[string(msg.Target)]
				redisMessages[string(msg.Target)] = id + 1
				redis := utils.GetRedis()
				var message = models.Message{
					ID:       id + 1,
					Target:   msg.Target,
					Type:     msg.Type,
					UserName: msg.UserName,
					Uid:      msg.Uid,
					Message:  msg.Message,
				}
				marshal, _ := json.Marshal(message)
				redis.Set(context.Background(), string(msg.Target)+"_"+strconv.Itoa(int(id+1)), marshal, 0)

				msg.ID = id + 1
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
		if len(messages) == 0 {
			redisMessages[cid] = 0
		} else {
			redisMessages[cid] = messages[0].ID
		}

		for _, message := range messages {
			marshal, _ := json.Marshal(message)
			err := redis.Set(context.Background(), cid+"_"+strconv.Itoa(int(message.ID)), marshal, 0).Err()
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

	var cids []string
	result := db.Model(&models.ChatRoom{}).Select("cid").Find(&cids)
	if result.Error != nil {
		panic(result.Error)
	}
	for _, cid := range cids {
		var latest uint
		result = db.Table("messages_" + cid).Model(&models.Message{}).Select("id").Order("id desc").First(&latest)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			latest = 0
		}
		keys, _ := redis.Keys(ctx, cid+"_*").Result()
		var neededInsertMessages []models.Message
		for _, key := range keys {
			split := strings.Split(key, "_")
			if id, _ := strconv.Atoi(split[1]); uint(id) > latest {
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
				redis.Del(ctx, key)
			}
		}

		if len(neededInsertMessages) != 0 {
			tx := db.Begin()
			result = tx.Table("messages_" + cid).Create(&neededInsertMessages)
			if result.Error != nil {
				tx.Rollback()
				panic(tx.Error)
			}
			tx.Commit()
		}
	}
}
