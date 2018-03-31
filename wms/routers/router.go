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
	beego.Router("/Post", &controllers.MessageController{}, "post:Post")
	beego.Router("/GetToken", &controllers.WechatController{}, "get:GetToken")
	beego.Router("/PostCustomMsg", &controllers.WechatController{}, "post:PostCustomMsg")
	beego.Router("/PostTemplateMsg", &controllers.WechatController{}, "post:PostTemplateMsg")
}
