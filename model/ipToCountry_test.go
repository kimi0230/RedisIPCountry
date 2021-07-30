package model

import (
	"RedisIPCountry/config"
	"RedisIPCountry/connect"
	"RedisIPCountry/utils"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestIpToScore(t *testing.T) {
	conn := connect.ConnectRedis()
	client := NewClient(conn)
	var tests = []struct {
		arg1 string
		want int64
	}{
		{
			"127.0.0.1",
			2130706433, // 0*256+127 -> 127*256+0 -> 32512*256+0 -> 8323072*256+1
		},
		{
			"111.71.213.250",
			1866978810, // 0*256+127 -> 127*256+0 -> 32512*256+0 -> 8323072*256+1
		},
	}

	for _, tt := range tests {
		if got := client.IpToScore(tt.arg1); got != tt.want {
			t.Errorf("got = %v, want = %v", got, tt.want)
		}
	}
}

func TestImportIpsToRedis(t *testing.T) {
	var ctx = context.Background()
	conn := connect.ConnectRedis()
	client := NewClient(conn)
	defer client.Conn.FlushDB(ctx)

	client.ImportIpsToRedis(config.FilePath + "GeoLite2-City-Blocks-IPv4.csv")
}

func TestImportCityToRedis(t *testing.T) {
	var ctx = context.Background()
	conn := connect.ConnectRedis()
	client := NewClient(conn)
	defer client.Conn.FlushDB(ctx)

	client.ImportCityToRedis(config.FilePath + "GeoLite2-City-Locations-en.csv")
}

//  go test  -run ^TestALL$ RedisIPCountry/model -v
func TestALL(t *testing.T) {
	var ctx = context.Background()
	conn := connect.ConnectRedis()
	client := NewClient(conn)

	t.Run("Test ip lookup", func(t *testing.T) {
		t.Log("Importing IP addresses to Redis...")
		client.ImportIpsToRedis(config.FilePath + "GeoLite2-City-Blocks-IPv4.csv")
		ranges := client.Conn.ZCard(ctx, "ip2cityid:").Val()
		t.Log("Loaded ranges into Redis:", ranges)
		utils.AssertTrue(t, ranges > 1000)

		t.Log("Importing Location lookups to Redis...")
		client.ImportCityToRedis(config.FilePath + "GeoLite2-City-Locations-en.csv")
		cities := client.Conn.HLen(ctx, "cityid2city:").Val()
		t.Log("Loaded city lookups into Redis:", cities)
		utils.AssertTrue(t, cities > 1000)

		// 隨機測試ip
		for i := 0; i < 5; i++ {
			ip := fmt.Sprintf("%s.%s.%s.%s", strconv.Itoa(rand.Intn(254)+1), utils.RandomString(256), utils.RandomString(256), utils.RandomString(256))
			t.Log(ip, client.FindCityByIp(ip))
		}

		defer client.Conn.FlushDB(ctx)
	})
}
