package initialize

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"online_Shop/userop_srv/global"
)

func InitRedisLock() {
	c := global.ServerConfig.RedisConfigInfo
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", c.Host, c.Port),
	})
	pool := goredis.NewPool(client)

	global.RDBLock = redsync.New(pool)

}
