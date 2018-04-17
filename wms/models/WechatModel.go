package models

import (
	"errors"
	"fmt"
	"message_server/utils/curl"
	"message_server/utils/db"
	"message_server/utils/logs"
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

var (
	prefix_secret   = "WXSCT_"
	prefix_template = "WXTPL_"
)

func GetToken(appid, secret string) (string, error) {
	var token string
	if appid == "" {
		return "", errors.New("No Appid")
	}
	if secret == "" {
		return "", errors.New("No Secret")
	}
	// redis
	redis_conn := db.RedisPool.Get()
	defer redis_conn.Close()
	exists, err := db.Bool(redis_conn.Do("EXISTS", appid))
	if err != nil {
		return "", err
	}
	if exists {
		token, err = db.String(redis_conn.Do("GET", appid))
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

func GetSecret(appid string) (secret string, err error) {
	redis_conn := db.RedisPool.Get()
	defer redis_conn.Close()

	secret, err = db.String(redis_conn.Do("GET", prefix_secret+appid))
	return
}

func GetTemplate(appid, name string) (tplid string, err error) {
	redis_conn := db.RedisPool.Get()
	defer redis_conn.Close()

	tplid, err = db.String(redis_conn.Do("GET", prefix_template+appid+"_"+name))
	return
}

func InitWechat() {
	initSecret()
	initTemplate()
}

func InitToken() {

}

// 初始化密钥
func initSecret() {
	rows, err := db.MysqlConn.Query("SELECT appid, secret FROM m_business_app")
	if err != nil {
		logs.Error(err)
	}

	redis_conn := db.RedisPool.Get()
	defer redis_conn.Close()
	for rows.Next() {
		var appid, secret string
		err = rows.Scan(&appid, &secret)
		if err != nil {
			logs.Error(err)
		}
		redis_conn.Do("SET", prefix_secret+appid, secret)
	}
}

// 初始化模版
func initTemplate() {
	rows, err := db.MysqlConn.Query("SELECT a.appid,a.template_id AS tplid,b.name FROM `m_wxmsg_tmpl` AS a LEFT JOIN `m_wxmsg_type` AS b ON a.type_id = b.id")
	if err != nil {
		logs.Error(err)
	}

	redis_conn := db.RedisPool.Get()
	defer redis_conn.Close()
	for rows.Next() {
		var appid, tplid, name string
		err = rows.Scan(&appid, &tplid, &name)
		if err != nil {
			logs.Error(err)
		}
		redis_conn.Do("SET", prefix_template+appid+"_"+name, tplid)
	}
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
