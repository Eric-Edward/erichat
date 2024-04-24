package global

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
)

func InitGlobalGoroutines() {
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
				tableName := "messages_" + string(msg.Target)
				db := utils.GetMySQLDB()
				tx := db.Begin()
				err := tx.Table(tableName).AutoMigrate(models.Message{})
				if err != nil {
					tx.Rollback()
					panic(err)
				}
				var message = models.Message{
					Target:   msg.Target,
					Type:     msg.Type,
					Message:  msg.Message,
					UserName: msg.UserName,
					Uid:      msg.Uid,
				}
				tx.Table(tableName).Create(&message)
				if tx.Error != nil {
					tx.Rollback()
					panic(err)
				}
				tx.Commit()

				msg.ID = message.ID
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
			}
		}
	}()
}
