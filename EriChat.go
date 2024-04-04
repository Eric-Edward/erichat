package main

import (
	"EriChat/router"
	"EriChat/utils"
)

func main() {
	ginServer := router.Route()
	utils.InitConfig()

	err := ginServer.Run(":8081")
	if err != nil {
		return
	}
}
