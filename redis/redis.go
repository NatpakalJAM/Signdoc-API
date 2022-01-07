package redis

import (
	"context"
	"fmt"
	"signdoc_api/config"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

//Init init redis package
func Init() {
	rc := config.C.Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", rc.Host, rc.Port),
		Password: rc.Password,
		DB:       rc.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("Error connect Redis: %s", err))
	}

	fmt.Println("redis init completed.")
}
