package main

import (
	_ "message_server/routers"

	"github.com/astaxie/beego"
)

func main() {
	// if beego.BConfig.RunMode == "dev" {
	// 	beego.BConfig.WebConfig.DirectoryIndex = true
	// 	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }
	// m.InitToken()
	beego.Run()
}
