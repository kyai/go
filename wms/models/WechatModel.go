package models

import (
	"errors"
	"fmt"
	"message_server/utils/curl"
	"message_server/utils/logs"
	"message_server/utils/redis"
	"strings"

	"github.com/astaxie/beego"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	expires     int64
}

func GetToken(appid, secret string) (string, error) {
	var token string
	if appid == "" {
		return "", errors.New("No Appid")
	}
	if secret == "" {
		return "", errors.New("No Secret")
	}
	// redis
	redis_conn := redis.RedisPool.Get()
	defer redis_conn.Close()
	exists, err := redis.Bool(redis_conn.Do("EXISTS", appid))
	if err != nil {
		return "", err
	}
	if exists {
		token, err = redis.String(redis_conn.Do("GET", appid))
	} else {
		tmp, err := GetTokenWX(appid, secret)
		if err != nil {
			return "", err
		}
		token = tmp.AccessToken
		redis_conn.Do("SET", appid, token, "EX", tmp.ExpiresIn-1e2)
	}

	return token, err
}

func GetTokenWX(appid, secret string) (Token, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, secret)

	var t Token
	curl.Do("GET", url, nil, &t)
	if t.Errcode == 0 && t.AccessToken != "" {
		return t, nil
	} else {
		return t, errors.New(t.Errmsg)
	}
}

func InitToken() {

}

func GetTemplateId(m string, i int) string {
	tpl := beego.AppConfig.String("message::" + m)
	tpls := strings.Split(tpl, "|")
	if i < len(tpls) {
		return tpls[i]
	} else {
		return ""
	}
}

func PostMessage(url string, data map[string]interface{}) interface{} {
	var resp interface{}
	curl.Do("POST", url, data, &resp)
	logs.Debug(resp)
	return resp
}
