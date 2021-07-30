package model

import (
	"RedisIPCountry/utils"
	"context"
	"encoding/json"
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

/**
 * @description: 產生 zset: key = ip2cityid, member = $cityID, score = $resIP
 * @param {string} filename
 * @return {*}
 */
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

type cityInfo struct {
	CityId  string
	Country string
	Region  string
	City    string
}

/**
 * @description: 產生Hash key=cityid2city, field= $city.CityId, value=$value
 * @param {string} filename
 * @return {*}
 */
func (c *Client) ImportCityToRedis(filename string) {
	res := utils.CSVReader(filename)
	/*
		geoname_id,locale_code,continent_code,continent_name,country_iso_code,country_name,subdivision_1_iso_code,subdivision_1_name,subdivision_2_iso_code,subdivision_2_name,city_name,metro_code,time_zone,is_in_european_union
		1665148,en,AS,Asia,TW,Taiwan,NWT,"New Taipei",,,"New Taipei",,Asia/Taipei,0
	*/
	var ctx = context.Background()
	pipe := c.Conn.Pipeline()
	for count, row := range res {
		if len(row) < 4 || !utils.IsDigital(row[0]) {
			continue
		}

		city := cityInfo{
			CityId:  row[0],
			Region:  row[3],
			Country: row[5],
			City:    row[6],
		}
		value, err := json.Marshal(city)

		pipe.HSet(ctx, "cityid2city:", city.CityId, value)
		if err != nil {
			log.Println("marshal json failed, err: ", err)
		}
		if (count+1)%1000 == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				log.Println("pipeline err in ImportCityToRedis: ", err)
				return
			}
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("pipeline err in ImportCityToRedis: ", err)
		return
	}
}

/**
 * @description:
 * @param {string} ip
 * @return {*}
 */
func (c *Client) FindCityByIp(ip string) string {
	var ctx = context.Background()
	ipAddress := strconv.Itoa(int(c.IpToScore(ip)))
	// Min:最小分數, Max:"10" 最大分數, Offset:0 類似 sql 的 limit, Count: 一次返回多少數據
	res := c.Conn.ZRevRangeByScore(ctx, "ip2cityid:", &redis.ZRangeBy{Max: ipAddress, Min: "0", Offset: 0, Count: 2}).Val()
	if len(res) == 0 {
		return ""
	}

	// 從 ip2cityid (zset) 取出 city id 並在 cityid2city(hash)中找城市資訊
	cityId := strings.Split(res[0], "_")[0]
	var result cityInfo
	if err := json.Unmarshal([]byte(c.Conn.HGet(ctx, "cityid2city:", cityId).Val()), &result); err != nil {
		log.Fatalln("unmarshal err: ", err)
	}
	return strings.Join([]string{result.CityId, result.City, result.Country, result.Region}, " ")
}
