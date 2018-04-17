package db

import (
	"database/sql"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

var MysqlConn *sql.DB

func InitMysqlConn() {
	var err error
	mysql_server := beego.AppConfig.String("mysql::mysql_server")
	MysqlConn, err = sql.Open("mysql", mysql_server)
	if err != nil {
		panic(err)
	}

	err = MysqlConn.Ping()
	if err != nil {
		panic(err)
	}
	return
}
