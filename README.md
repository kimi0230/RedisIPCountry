# Redis IP to Country
Mapping IP to Country by [go-redis/redis/v8](https://github.com/go-redis/redis)

## Run
### config
`config/config.go`
``` go
// config for redis
var (
	// Addr = "redis-ipcountry-redis:6379" // for docker
	Addr     = "127.0.0.1:6379" // for localhost
	Password = ""
	DB       = 1
)

// file path for ip to country
const FilePath = "../asset/GeoLite2-City-CSV/"
```

### start service
``` shell
docker-compose up -d
```

### shutdown
``` shell
docker-compose down
```
