package logs

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func init() {
	// display file and filename
	logs.EnableFuncCallDepth(true)
	// level
	logs.SetLogFuncCallDepth(2)

	logs.Async(1e3)

	logs.SetLogger(logs.AdapterFile, `{"filename":"kms.log"}`)

	if beego.BConfig.RunMode == "pro" {
		logs.SetLevel(logs.LevelError)
	}
}

func Error(f interface{}, v ...interface{}) {
	logs.Error(f, v)
}

func Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v)
}

func Info(f interface{}) {
	logs.Informational(f)
}
