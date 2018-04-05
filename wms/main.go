package main

import (
	m "message_server/models"
	_ "message_server/routers"
	"message_server/utils/redis"

	"github.com/astaxie/beego"
)

func main() {
	// if beego.BConfig.RunMode == "dev" {
	// 	beego.BConfig.WebConfig.DirectoryIndex = true
	// 	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }

	// redis
	redis.InitRedisPool()

	// socket
	go m.InitSocket()

	beego.Run()
}
