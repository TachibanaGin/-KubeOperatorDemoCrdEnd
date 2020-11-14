package main

import (
	"Crd-End/mysql"
	"Crd-End/redis"
	"Crd-End/router"
)

func init()  {
	redis.InitRedis()
	mysql.InitMysql()
}
func main() {
	r := router.InitRouter()
	defer redis.Client.Close()
	defer mysql.Db.Close()
	r.Run(":7777") // 监听并在 0.0.0.0:7777 上启动服务
}