package main

import (
	m "message_server/models"
	_ "message_server/routers"
	"message_server/utils/db"
	_ "message_server/utils/proc"

	"github.com/astaxie/beego"
)

func main() {
	// if beego.BConfig.RunMode == "dev" {
	// 	beego.BConfig.WebConfig.DirectoryIndex = true
	// 	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }

	// mysql
	db.InitMysqlConn()
	defer db.MysqlConn.Close()

	// redis
	db.InitRedisPool()

	m.InitWechat()

	// socket
	go m.InitSocket()

	beego.Run()
}
