package controllers

import "github.com/astaxie/beego"

type DefaultController struct {
	beego.Controller
}

func (c *DefaultController) Get() {
	c.Ctx.WriteString("Hello Golang!")
	return
}
