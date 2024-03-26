package main

import (
	"EcChat/router"
	"EcChat/utils"
)

func main() {
	ginServer := router.Route()
	utils.InitConfig()

	err := ginServer.Run(":8081")
	if err != nil {
		return
	}
}
