package config

// config for redis
var (
	// Addr = "redis-ipcountry-redis:6379" // for docker
	Addr     = "127.0.0.1:6379" // for localhost
	Password = ""
	DB       = 1
)

// file path for ip to country
const FilePath = "../asset/GeoLite2-City-CSV/"
