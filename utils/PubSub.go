package utils

import (
	"context"
	"fmt"
)

//TODO 这里的pubsub之后需要改成缓存就可以了

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
