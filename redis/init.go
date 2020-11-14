package redis

import (
	"Crd-End/config"
	"fmt"
	"github.com/go-redis/redis/v7"
)

var(
	redisConf = config.LoadRedis()
	host = redisConf["host"]
	port = redisConf["port"]
	password = redisConf["password"]
	Client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password, // no password set
		//DB:       0,  // use default DB
	})
)

func InitRedis(){

	pong, err := Client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}
