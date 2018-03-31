package controllers

import (
	"encoding/json"
	"strings"

	m "message_server/models"
	"message_server/utils/curl"
	"message_server/utils/logs"

	"github.com/astaxie/beego"
)

type MessageController struct {
	BaseController
}

type Request struct {
	Message string      `json:"message"`
	Tousers []string    `json:"tousers"`
	Link    string      `json:"link"`
	Data    interface{} `json:"data"`
}

func (c *MessageController) Post() {
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

	tpl := beego.AppConfig.String("message::" + r.Message)
	tpl_wxgzh := strings.Split(tpl, "|")[0]
	// tpl_wxapp := strings.Split("|", tpl)[1]

	token, _ := m.GetToken()

	if tpl_wxgzh != "" {
		url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + token
		for _, v := range r.Tousers {
			data := make(map[string]interface{})
			data["touser"] = v
			data["template_id"] = tpl_wxgzh
			data["url"] = r.Link
			data["data"] = r.Data
			curl.Do("POST", url, data, nil)
		}
	}
}
