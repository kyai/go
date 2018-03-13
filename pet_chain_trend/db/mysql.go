package db

import (
	"database/sql"

	"github.com/dlintw/goconf"
	_ "github.com/go-sql-driver/mysql"
)

var MysqlConn *sql.DB

func InitMysqlConn(conf *goconf.ConfigFile) {
	var err error
	mysql_server, _ := conf.GetString("mysql", "mysql_server")
	MysqlConn, err = sql.Open("mysql", mysql_server)
	if err != nil {
		panic(err)
	}
	return
}
