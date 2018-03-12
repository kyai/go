package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var MysqlConn *sql.DB

func InitMysqlConn() {
	var err error
	MysqlConn, err = sql.Open("mysql", "root:12345678@tcp(localhost)/pet_chain_trend?charset=utf8")
	if err != nil {
		panic(err)
	}
	return
}
