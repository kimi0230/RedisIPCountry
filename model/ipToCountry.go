package model

import (
	"RedisIPCountry/utils"
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	Conn *redis.Client
}

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}

// 將IP轉成整數之後, 程式就可以創建IP address與 城市 id 之間的mapping
func (c *Client) IpToScore(ip string) int64 {
	var score int64 = 0
	for _, v := range strings.Split(ip, ".") {
		n, _ := strconv.ParseInt(v, 10, 0)
		score = score*256 + n
	}
	return score
}

func (c *Client) ImportIpsToRedis(filename string) {
	res := utils.CSVReader(filename)
	/*
		network,geoname_id,registered_country_geoname_id,represented_country_geoname_id,is_anonymous_proxy,is_satellite_provider,postal_code,latitude,longitude,accuracy_radius
		1.0.0.0/24,2077456,2077456,,0,0,,-33.4940,143.2104,1000
	*/
	var ctx = context.Background()
	pipe := c.Conn.Pipeline()
	for count, row := range res {
		var (
			startIp string
			resIP   int64
		)
		if len(row) == 0 {
			startIp = ""
		} else {
			startIp = row[0]
		}

		if strings.Contains(strings.ToLower(startIp), "i") {
			continue
		}
		if strings.Contains(startIp, ".") {
			resIP = c.IpToScore(startIp)
		} else {
			var err error
			resIP, err = strconv.ParseInt(startIp, 10, 64)
			if err != nil {
				continue
			}
		}
		// 因為多個IP範圍可能會對應到同一個城市 ID, 故在後面_ 加上目前已有城市ID數量
		cityID := row[2] + "_" + strconv.Itoa(count)
		pipe.ZAdd(ctx, "ip2cityid:", &redis.Z{Member: cityID, Score: float64(resIP)})
		if (count+1)%1000 == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				log.Println("pipeline err in ImportIpsToRedis: ", err)
				return
			}
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("pipeline err in ImportIpsToRedis: ", err)
		return
	}
}
