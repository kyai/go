package controllers

import (
	"encoding/json"

	m "message_server/models"
	"message_server/utils/logs"

	"github.com/astaxie/beego"
)

type WechatDevController struct {
	BaseController
}

func (c *WechatDevController) GetTokenTest() {
	appid := beego.AppConfig.String("wechat_test::appid")
	secret := beego.AppConfig.String("wechat_test::secret")
	token, err := m.GetToken(appid, secret)
	if err != nil {
		c.JSON("102", "error", err.Error())
	} else {
		c.JSON("101", "success", token)
	}
}

func (c *WechatDevController) PostTest() {
	var r Request
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &r); err != nil {
		logs.Error(err.Error())
		c.JSON("102", "error", err.Error())
		return
	}

	if r.Message == "" {
		c.JSON("102", "error", "Undefined message type")
		return
	}

	appid := beego.AppConfig.String("wechat_test::appid")
	secret := beego.AppConfig.String("wechat_test::secret")
	token, err := m.GetToken(appid, secret)
	if err != nil {
		c.JSON("102", "error", err.Error())
		return
	}

	if tplid := beego.AppConfig.String("message_test::" + r.Message); tplid != "" {
		url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + token
		for _, v := range r.Tousers {
			data := make(map[string]interface{})
			data["touser"] = v
			data["template_id"] = tplid
			data["url"] = r.Link
			data["data"] = r.Data
			go m.PostMessage(url, data)
		}
	} else {
		logs.Error("Has not template id", r.Message)
	}
}
