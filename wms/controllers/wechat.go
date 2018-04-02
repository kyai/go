package controllers

import (
	"encoding/json"

	m "message_server/models"
	"message_server/utils/logs"

	"github.com/astaxie/beego"
)

type WechatController struct {
	BaseController
}

type Request struct {
	Message string      `json:"message"`
	Tousers []string    `json:"tousers"`
	Link    string      `json:"link"`
	Data    interface{} `json:"data"`
}

func (c *WechatController) GetToken() {
	appid := beego.AppConfig.String("wechat::appid")
	secret := beego.AppConfig.String("wechat::secret")
	token, err := m.GetToken(appid, secret)
	if err != nil {
		c.JSON("102", "error", err.Error())
	} else {
		c.JSON("101", "success", token)
	}
}

func (c *WechatController) GetTokenApp() {
	appid := c.GetString("appid")
	secret := c.GetString("secret")
	token, err := m.GetToken(appid, secret)
	if err != nil {
		c.JSON("102", "error", err.Error())
	} else {
		c.JSON("101", "success", token)
	}
}

func (c *WechatController) Post() {
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

	appid := beego.AppConfig.String("wechat::appid")
	secret := beego.AppConfig.String("wechat::secret")
	token, err := m.GetToken(appid, secret)
	if err != nil {
		c.JSON("102", "error", err.Error())
		return
	}

	if tplid := m.GetTemplateId(r.Message, 0); tplid != "" {
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

// 发送到小程序
func (c *WechatController) PostToProgram() {
	var r struct {
		Message string      `json:"message"`
		Tousers []string    `json:"tousers"`
		Form    string      `json:"form"`
		Link    string      `json:"link"`
		Data    interface{} `json:"data"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &r); err != nil {
		logs.Error(err.Error())
		c.JSON("102", "error", err.Error())
		return
	}

	if r.Message == "" {
		c.JSON("102", "error", "Undefined message type")
		return
	}

	appid := c.GetString("appid")
	secret := c.GetString("secret")
	token, err := m.GetToken(appid, secret)
	if err != nil {
		c.JSON("102", "error", err.Error())
		return
	}

	if tplid := m.GetTemplateId(r.Message, 1); tplid != "" {
		url := "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=" + token
		for _, v := range r.Tousers {
			data := make(map[string]interface{})
			data["touser"] = v
			data["template_id"] = tplid
			data["form_id"] = r.Form
			data["page"] = r.Link
			data["data"] = r.Data
			go m.PostMessage(url, data)
		}
	} else {
		logs.Error("Has not template id", r.Message)
	}
}
