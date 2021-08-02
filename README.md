# Redis IP to Country
Mapping IP to Country use [go-redis/redis/v8](https://github.com/go-redis/redis)

## GeoLite2 City
1. sign up : https://www.maxmind.com/en/home
2. guide : https://blog.csdn.net/qq_26373925/article/details/111876765
## Run
### config
Path : `config/config.go`

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

### Reference
* [go-redis 連接池](https://www.huaweicloud.com/articles/db24f1e8b4a4f0218ddf08463d8ec871.html)