package global

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
)

var persistenceData chan utils.WsMessage

func InitGlobalGoroutines() {
	persistenceData = make(chan utils.WsMessage)
	persistenceDataEventLoop()

}

func persistenceDataEventLoop() {
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
					SendBy:  string(msg.Uid),
					Content: msg.Message,
					Url:     "",
					Size:    int64(len(msg.Message)),
				})
				if tx.Error != nil {
					tx.Rollback()
					panic(err)
				}
				tx.Commit()
			}
		}
	}()
}
func PersistenceData() chan utils.WsMessage {
	return persistenceData
}
