package models

import (
	"github.com/go-redis/redis"
	"log"
)

var client *redis.Client

//初始化数据库链接
func init(){
	client=redis.NewClient(&redis.Options{
		Addr:     "172.18.0.2:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	if err!=nil{
		log.Fatal("redis 链接异常")
	}
	log.Println("redis connected and response:", pong)
}
