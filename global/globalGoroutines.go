package global

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
)

var persistenceData chan utils.WsMessage
var confirmData chan utils.WsMessage

func InitGlobalGoroutines() {
	persistenceData = make(chan utils.WsMessage)
	confirmData = make(chan utils.WsMessage)
	messageEventLoop()
}

func messageEventLoop() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			select {
			case msg := <-persistenceData:
				tableName := "messages_" + string(msg.Target)
				db := utils.GetMySQLDB()
				tx := db.Begin()
				err := tx.Table(tableName).AutoMigrate(models.Message{})
				if err != nil {
					tx.Rollback()
					panic(err)
				}
				tx.Table(tableName).Create(&models.Message{
					Target:   msg.Target,
					Type:     msg.Type,
					Message:  msg.Message,
					UserName: msg.UserName,
					Uid:      msg.Uid,
				})
				if tx.Error != nil {
					tx.Rollback()
					panic(err)
				}
				tx.Commit()
			case msg := <-confirmData:
				db := utils.GetMySQLDB()
				tx := db.Begin()
				result := tx.Model(&models.ChatRoomMember{}).Where("cid=? and uid=? and record<?", msg.Target, msg.Uid, msg.Message).Update("record", msg.Message)
				if result.Error != nil {
					tx.Rollback()
					panic(tx.Error)
				}
			}
		}
	}()
}
func PersistenceData() chan utils.WsMessage {
	return persistenceData
}
func ConfirmData() chan utils.WsMessage { return confirmData }
