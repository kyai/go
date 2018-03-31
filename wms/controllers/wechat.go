package controllers

import (
	"encoding/json"
	"fmt"
	m "message_server/models"
)

type WechatController struct {
	BaseController
}

func init() {
}

func (c *WechatController) GetToken() {
	token, err := m.GetToken()
	if err != nil {
		c.JSON("102", "error", err.Error())
	} else {
		c.JSON("101", "success", token)
	}
}

type ReqCustomMsg struct {
	Tousers []string `json:"tousers"`
	Content string   `json:"content"`
}

func (c *WechatController) PostCustomMsg() {
	var req ReqCustomMsg
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		fmt.Println(err.Error())
		return
	}

	token, _ := m.GetToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + token

	for _, v := range req.Tousers {
		data := make(map[string]interface{})
		data["touser"] = v
		data["msgtype"] = "text"
		text := make(map[string]interface{})
		text["content"] = req.Content
		data["text"] = text
		c.CURL("POST", url, data, nil)
	}
}

type ReqTemplateMsg struct {
	Tousers    []string    `json:"tousers"`
	TemplateId string      `json:"template_id"`
	Url        string      `json:"url"`
	Data       interface{} `json:"data"`
}

func (c *WechatController) PostTemplateMsg() {
	var req ReqTemplateMsg
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		fmt.Println(err.Error())
		c.JSON("102", "error", err.Error())
		return
	}

	token, _ := m.GetToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + token

	for _, v := range req.Tousers {
		data := make(map[string]interface{})
		data["touser"] = v
		data["template_id"] = req.TemplateId
		data["url"] = req.Url
		data["data"] = req.Data
		c.CURL("POST", url, data, nil)
	}
}
