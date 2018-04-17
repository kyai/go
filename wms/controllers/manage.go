package controllers

import (
	m "message_server/models"
	"message_server/utils/file"
	"time"
)

type ManageController struct {
	BaseController
}

func (c *ManageController) GetLog() {
	date := c.GetString("date")
	if date != "" {
		timep, err := time.Parse("2006-01-02", date)
		if err != nil {
			c.JSON("102", "error", err.Error())
		}
		date = timep.Format("2006-01-02")
		if date == time.Now().Format("2006-01-02") {
			date = ""
		} else {
			date += "."
		}
	}
	path := "logs/app." + date + "log"

	data, err := file.Reader(path)
	if err != nil {
		c.JSON("102", "error", err.Error())
	} else {
		c.JSON("101", "success", data)
	}
}

func (c *ManageController) ReWxCache() {
	m.InitWechat()
	c.JSON("101", "success", nil)
}
