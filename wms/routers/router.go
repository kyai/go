package routers

import (
	"message_server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// ns := beego.NewNamespace("/v1",
	// 	beego.NSNamespace("/default",
	// 		beego.NSInclude(
	// 			&controllers.DefaultController{},
	// 		),
	// 	),
	// )
	// beego.AddNamespace(ns)

	beego.Router("/", &controllers.DefaultController{})
	beego.Router("/Post", &controllers.WechatController{}, "post:Post")
	beego.Router("/PostToProgram", &controllers.WechatController{}, "post:PostToProgram")
	beego.Router("/GetToken", &controllers.WechatController{}, "get:GetToken")
	beego.Router("/GetTokenApp", &controllers.WechatController{}, "get:GetTokenApp")

	// test
	beego.Router("/GetTokenTest", &controllers.WechatDevController{}, "get:GetTokenTest")
	beego.Router("/PostTest", &controllers.WechatDevController{}, "post:PostTest")
	beego.Router("/Send", &controllers.SocketController{}, "get:Send")
	beego.Router("/GetSecret", &controllers.WechatController{}, "get:GetSecret")
	beego.Router("/GetTemplate", &controllers.WechatController{}, "get:GetTemplate")

	// manage
	beego.Router("/manage_api/GetLog", &controllers.ManageController{}, "get:GetLog")
	beego.Router("/manage_api/ReWxCache", &controllers.ManageController{}, "get:ReWxCache")
}
