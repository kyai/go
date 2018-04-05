package controllers

import (
	m "message_server/models"
)

type SocketController struct {
	BaseController
}

func (c *SocketController) Send() {
	uid := c.GetString("uid")
	msg := c.GetString("msg")
	if err := m.Send(uid, msg); err != nil {
		c.JSON("102", "error", err.Error())
	} else {
		c.JSON("101", "success", nil)
	}
}
