package controllers

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type DateController struct {
	beego.Controller
}

// @Title Get
// @Description get date by format
// @Param	format		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /date [get]
func (c *DateController) Get() {
	format := c.GetString("format")
	if format == "" {
		format = "yyyy-MM-dd hh:mm:ss"
	}

	re := func(old, new string) {
		format = strings.Replace(format, old, new, 1)
	}

	re("yyyy", "2006")
	re("yy", "06")
	re("MM", "01")
	re("M", "1")
	re("dd", "02")
	re("d", "2")
	re("H", "15")
	re("hh", "03")
	re("h", "3")
	re("mm", "04")
	re("m", "4")
	re("ss", "05")
	re("s", "5")

	format = time.Now().Format(format)
	c.Ctx.WriteString(format)
	return
}
