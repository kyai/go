package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pet_chain_trend/db"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	jsonR   map[string]interface{}
	httpC   http.Client
	request *http.Request
)

func init() {
	jsonR = make(map[string]interface{})
	jsonR["pageNo"] = 1
	jsonR["pageSize"] = 20
	jsonR["querySortType"] = "AMOUNT_ASC"
	jsonR["petIds"] = []int{}
	jsonR["requestId"] = 1

	jsonS, err := json.Marshal(jsonR)
	if err != nil {
		fmt.Println(err.Error())
	}
	reader := bytes.NewReader(jsonS)
	url := "https://pet-chain.baidu.com/data/market/queryPetsOnSale"
	request, err = http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
	}

	request.Header.Set("Content-Type", "application/json")

	httpC = http.Client{}
}

// The structure of the response data
type Market struct {
	Data struct {
		PetsOnSale []struct {
			Amount float32 `json:"amount,string"`
		} `json:"petsOnSale"`
	} `json:"data"`
}

// Get data and save it to the database(redis)
// Run the collector
func Run() {
	now, avg, err := post()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	now_num, err := strconv.Atoi(now.Format("2006010215"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	save(now_num, avg)
}

func post() (time.Time, int, error) {
	now := time.Now()
	avg := 0

	resp, err := httpC.Do(request)
	if err != nil {
		return now, avg, err
	}

	rdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return now, avg, err
	}

	var market Market
	if err = json.Unmarshal(rdata, &market); err != nil {
		return now, avg, err
	}

	if arr := market.Data.PetsOnSale; len(arr) > 0 {
		var sum float32
		for _, v := range arr {
			sum += v.Amount
		}
		avg = int(sum / float32(len(arr)))
	}

	return now, avg, err
}

func save(now, avg int) {
	// redis
	conn := db.RedisPool.Get()
	defer conn.Close()

	pre, err := redis.Int(conn.Do("GET", "PCT_TIME_K"))
	if err != nil {
		exists, err := redis.Bool(conn.Do("EXISTS", "PCT_TIME_K"))
		if err != nil {
			panic(err)
		}
		if !exists {
			pre = now
		}
	}

	if now < pre {
		return
	}

	if now > pre {
		values, _ := redis.Ints(conn.Do("lrange", "PCT_TIME_V", "0", "-1"))

		// reset key and value
		_, _ = conn.Do("SET", "PCT_TIME_K", now)
		_, _ = conn.Do("DEL", "PCT_TIME_V")

		var sum int
		for _, v := range values {
			sum += v
		}
		avg := sum / len(values)

		saveToMysql(now, avg)
		return
	}

	_, err1 := conn.Do("SET", "PCT_TIME_K", now)
	_, err2 := conn.Do("lpush", "PCT_TIME_V", avg)
	if err1 != nil || err2 != nil {
		fmt.Println(err1.Error())
		return
	}
}

func saveToMysql(now, avg int) {
	_, err := db.MysqlConn.Exec("INSERT INTO record (time, average) VALUES (?, ?)", now, avg)
	if err != nil {
		panic(err)
	}
}
