package model

import (
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
