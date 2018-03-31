package controllers

type DefaultController struct {
	BaseController
}

func (c *DefaultController) Get() {
	c.Ctx.WriteString("Hello KTV")
}
