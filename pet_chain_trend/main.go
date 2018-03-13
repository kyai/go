package main

import (
	"flag"
	"fmt"
	"pet_chain_trend/db"
	"pet_chain_trend/server"
	"time"

	"github.com/dlintw/goconf"
)

func main() {

	// read config
	conf_file := flag.String("config", "./config.ini", "set config file.")
	flag.Parse()
	conf, err := goconf.ReadConfigFile(*conf_file)
	if err != nil {
		panic(err)
	}

	// redis
	db.InitRedisPool(conf)
	defer db.RedisPool.Close()
	// mysql
	db.InitMysqlConn(conf)
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
