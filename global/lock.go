package global

import (
	"EriChat/models"
	"EriChat/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"math/rand"
	"sort"
	"time"
)

// 重试次数
var retryTimes = 5

// 重试频率
var retryInterval = time.Millisecond * 50

var rdb *redis.Client

// 锁的默认过期时间
var expiration time.Duration

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Tsinghua",
		DB:       0,
	})
}

// MockTest 模拟分布式业务加锁场景
func MockTest(tag, cid string) {
	var ctx, cancel = context.WithCancel(context.Background())

	// 随机value
	lockV := getRandValue()

	lockK := "REDIS_LOCK"

	// 默认过期时间
	expiration = time.Millisecond * 200

	set, err := rdb.SetNX(ctx, lockK, lockV, expiration).Result()

	if err != nil {
		panic(err.Error())
	}
	if set == false && retry(ctx, rdb, lockK, lockV, expiration, tag) == false {
		fmt.Println(tag + " server unavailable, try again later")
		panic("重新加锁失败")
	}

	// 加锁成功,新增守护线程
	go watchDog(ctx, rdb, lockK, expiration, tag)

	// 处理业务(通过随机时间延迟模拟)
	var redisMessage RedisMessage
	bMessage, err := rdb.Get(ctx, cid).Result()
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal([]byte(bMessage), &redisMessage)
	if err != nil {
		panic(err.Error())
	}

	var latest uint
	db := utils.GetMySQLDB()
	result := db.Table("messages_" + cid).Model(&models.Message{}).Select("id").Order("id desc").First(&latest)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic(result.Error)
	}

	insertPos := sort.Search(len(redisMessage.Messages), func(i int) bool {
		return redisMessage.Messages[i].ID > latest
	})
	var insertMessage = redisMessage.Messages[insertPos:]
	if len(insertMessage) != 0 {
		tx := db.Begin()
		result = tx.Table("messages_" + cid).Model(&models.Message{}).Create(&insertMessage)
		if result.Error != nil {
			tx.Rollback()
			panic(result.Error)
		}
		tx.Commit()
	}

	if len(redisMessage.Messages) > 100 {
		redisMessage.Messages = redisMessage.Messages[len(redisMessage.Messages)-100 : len(redisMessage.Messages)]
	}

	marshal, err := json.Marshal(redisMessage)
	if err != nil {
		panic(err.Error())
	}
	err = rdb.Set(ctx, cid, marshal, 0).Err()
	if err != nil {
		panic(err.Error())
	}
	// 业务处理完成
	// 释放锁
	defer func() {
		_ = delByKeyWhenValueEquals(ctx, rdb, lockK, lockV)
		cancel()
	}()
}

// 释放锁
func delByKeyWhenValueEquals(ctx context.Context, rdb *redis.Client, key string, value interface{}) bool {
	lua := `
-- 如果当前值与锁值一致,删除key
if redis.call('GET', KEYS[1]) == ARGV[1] then
	return redis.call('DEL', KEYS[1])
else
	return 0
end
`
	scriptKeys := []string{key}

	val, err := rdb.Eval(ctx, lua, scriptKeys, value).Result()
	if err != nil {
		panic(err.Error())
	}

	return val == int64(1)
}

// 生成随机时间
func getRandDuration() time.Duration {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	minT := 50
	maxT := 100
	return time.Duration(rand.Intn(maxT-minT)+minT) * time.Millisecond
}

// 生成随机值
func getRandValue() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Int()
}

// 守护线程
func watchDog(ctx context.Context, rdb *redis.Client, key string, expiration time.Duration, tag string) {
	for {
		select {
		// 业务完成
		case <-ctx.Done():
			fmt.Printf("%s任务完成,关闭%s的自动续期\n", tag, key)
			return
			// 业务未完成
		default:
			// 自动续期
			rdb.PExpire(ctx, key, expiration)
			// 继续等待
			time.Sleep(expiration / 2)
		}
	}
}

// 重试
func retry(ctx context.Context, rdb *redis.Client, key string, value interface{}, expiration time.Duration, tag string) bool {
	i := 1
	for i <= retryTimes {
		fmt.Printf(tag+"第%d次尝试加锁中...\n", i)
		set, err := rdb.SetNX(ctx, key, value, expiration).Result()

		if err != nil {
			panic(err.Error())
		}

		if set == true {
			return true
		}

		time.Sleep(retryInterval)
		i++
	}
	return false
}
