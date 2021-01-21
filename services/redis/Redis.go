package redisClient

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
)

//直接连接
func Connect() redis.Conn{
	pool,_ := redis.Dial("tcp", beego.AppConfig.String("redisdb"))
	return pool
}

//连接池连接
func PoolConnect() redis.Conn{
	//建立连接池
	pool := &redis.Pool{
		MaxIdle:     5000,              //最大空闲连接数
		MaxActive:   10000,             //最大连接数
		IdleTimeout: 180 * time.Second, //空闲连接超时时间
		Wait:        true,              //超过最大连接数时，是等待还是报错
		Dial: func() (redis.Conn, error) { //建立链接
			c, err := redis.Dial("tcp", beego.AppConfig.String("redisdb"))
			if err != nil {
				return nil, err
			}
			// 选择db
			//c.Do("SELECT", '')
			return c, nil
		},
	}
	return pool.Get()
}
