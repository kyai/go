package controllers

import (
	"message_server/utils/curl"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) CURL(method, url string, data map[string]interface{}, obj interface{}) {
	curl.Do(method, url, data, obj)
}

func (c *BaseController) JSON(code, msg string, data interface{}) {
	var jdata = make(map[string]interface{})
	jdata["Code"] = code
	jdata["Msg"] = msg
	jdata["Data"] = data
	c.Data["json"] = jdata
	c.ServeJSON()
	return
}
