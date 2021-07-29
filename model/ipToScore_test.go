package model

import (
	"RedisIPCountry/connect"
	"testing"
)

var tests = []struct {
	arg1 string
	want int64
}{
	{
		"127.0.0.1",
		2130706433, // 0*256+127 -> 127*256+0 -> 32512*256+0 -> 8323072*256+1
	},
}

func TestIpToScore(t *testing.T) {
	conn := connect.ConnectRedis()
	client := NewClient(conn)
	for _, tt := range tests {
		if got := client.IpToScore(tt.arg1); got != tt.want {
			t.Errorf("got = %v, want = %v", got, tt.want)
		}
	}
}
