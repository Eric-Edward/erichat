package main

import (
	"EriChat/global"
	"EriChat/router"
	"EriChat/utils"
)

func main() {
	ginServer := router.Route()
	utils.InitConfig()
	global.InitGlobalGoroutines()

	err := ginServer.Run(":8081")
	if err != nil {
		return
	}
}
