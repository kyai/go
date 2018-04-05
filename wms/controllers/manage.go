package controllers

import "message_server/utils/file"

type ManageController struct {
	BaseController
}

func (c *ManageController) GetLog() {
	date := c.GetString("date")
	if date != "" {
		date = date + "."
	}
	path := "kms." + date + "log"

	data, err := file.Reader(path)
	if err != nil {
		c.JSON("102", "error", err.Error())
	} else {
		c.JSON("101", "success", data)
	}
}
