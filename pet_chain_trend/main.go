package main

import (
	"fmt"
	"pet_chain_trend/db"
	"pet_chain_trend/server"
	"time"
)

func main() {

	// redis
	db.InitRedisPool()
	defer db.RedisPool.Close()
	// mysql
	db.InitMysqlConn()
	defer db.MysqlConn.Close()
	// server
	ticker := time.NewTicker(time.Second * 10)
	func() {
		for k := range ticker.C {
			fmt.Println(k)
			server.Run()
		}
	}()

}
