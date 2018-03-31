package models

import (
	"errors"
	"fmt"
	"message_server/utils/curl"
	"message_server/utils/logs"
	"time"

	"github.com/astaxie/beego"
)

var AccessToken string
var appid, secret string

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

type MessageTpl struct {
	Touser     string      `json:"touser"`
	TemplateId string      `json:"template_id"`
	Url        string      `json:"url"`
	Data       interface{} `json:"data"`
}

func init() {
	appid = beego.AppConfig.String("wechat::appid")
	secret = beego.AppConfig.String("wechat::secret")
}

func GetToken() (string, error) {
	appid := beego.AppConfig.String("wechat::appid")
	secret := beego.AppConfig.String("wechat::secret")
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, secret)
	var token Token
	curl.Do("GET", url, nil, &token)
	if token.Errcode == 0 {
		return token.AccessToken, nil
	} else {
		logs.Error(token.Errmsg)
		return "", errors.New(token.Errmsg)
	}
}

func InitToken() {
	AccessToken, err := GetToken()
	if err != nil {
		logs.Error(err)
	}
	ticker := time.NewTicker(time.Second * time.Duration(5000))
	func() {
		for k := range ticker.C {
			AccessToken, err = GetToken()
			if err != nil {
				logs.Error(err)
			}
			fmt.Println(k)
			fmt.Println(AccessToken)
		}
	}()
}
