package lache

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhouhp1295/lache/driver"
	"time"
)

type DriverType int

const (
	Local        DriverType = 1
	Redis        DriverType = 2
	RedisCluster DriverType = 3
)

type Client struct {
	Driver Driver
}

func New(t DriverType, options any) *Client {
	client := new(Client)
	switch t {
	case Local:
		if opts, ok := options.(driver.LocalOptions); ok {
			client.Driver = driver.NewLocalDriver(opts)
		} else {
			panic("需要 LocalOptions")
		}
	case Redis:
		if opts, ok := options.(redis.Options); ok {
			client.Driver = driver.NewRedisDriver(opts)
		} else {
			panic("需要 redis.Options")
		}
	case RedisCluster:
		if opts, ok := options.(redis.ClusterOptions); ok {
			client.Driver = driver.NewRedisClusterDriver(opts)
		} else {
			panic("需要 redis.ClusterOptions")
		}
	default:
		panic("未知的缓存类型")
	}
	return client
}

func (client *Client) Get(key string) (result any, ok bool) {
	result, ok = client.Driver.Get(key)
	return
}

func (client *Client) GetT(key string, result any) (ok bool) {
	ok = client.Driver.GetT(key, result)
	return
}

func (client *Client) Set(key string, value any, expiration time.Duration) (ok bool) {
	ok = client.Driver.Set(key, value, expiration)
	return
}

func (client *Client) Delete(key string) (ok bool) {
	ok = client.Driver.Delete(key)
	return
}
