package utils

import (
	"context"
	"fmt"
)

func Publish(ctx context.Context, channel string, message string) {
	rdb := GetRedis()
	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		panic(err)
	}
}

func Subscribe(ctx context.Context, channel string) {
	rdb := GetRedis()
	sub := rdb.Subscribe(ctx, channel)
	defer func() {
		err := sub.Close()
		fmt.Println("订阅关闭失败！", err)
	}()

	messages := sub.Channel()
	for message := range messages {
		fmt.Println(message)
	}
}
