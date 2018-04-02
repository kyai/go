package routers

import (
	"message_server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/default",
			beego.NSInclude(
				&controllers.DefaultController{},
			),
		),
	)
	beego.AddNamespace(ns)

	beego.Router("/", &controllers.DefaultController{})
	beego.Router("/Post", &controllers.WechatController{}, "post:Post")
	beego.Router("/PostToProgram", &controllers.WechatController{}, "post:PostToProgram")
	beego.Router("/GetToken", &controllers.WechatController{}, "get:GetToken")
	beego.Router("/GetTokenApp", &controllers.WechatController{}, "get:GetTokenApp")
}
