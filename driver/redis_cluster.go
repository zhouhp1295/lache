// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package driver

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"time"
)

type RedisCluster struct {
	client    *redis.ClusterClient
	ctx       context.Context
	connected bool
}

func NewRedisClusterDriver(options redis.ClusterOptions) *RedisCluster {
	driver := new(RedisCluster)
	onConnected := options.OnConnect

	options.OnConnect = func(ctx context.Context, cn *redis.Conn) error {
		fmt.Printf("Redis Connected!\n")
		driver.connected = true
		if onConnected != nil {
			_ = onConnected(ctx, cn)
		}
		return nil
	}

	driver.client = redis.NewClusterClient(&options)

	driver.ctx = context.Background()

	err := driver.client.ClientGetName(driver.ctx).Err()
	if err != nil && err != redis.Nil {
		fmt.Printf("Redis Error = %s \n", err.Error())
	}

	return driver
}

func (d *RedisCluster) Get(key string) (result any, ok bool) {
	if !d.connected {
		return
	}
	var err error
	result, err = d.client.Get(d.ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("[RedisCluster Error][Get Key=%s] %s\n", key, err.Error())
		}
		ok = false
	} else {
		ok = true
	}
	return
}

func (d *RedisCluster) GetT(key string, result any) (ok bool) {
	if !d.connected {
		return
	}
	res, err := d.client.Get(d.ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("[Redis Error][GetT Key=%s] %s\n", key, err.Error())
		}
		return
	}

	ok = ParseString(res, result)
	return
}

func (d *RedisCluster) Set(key string, value any, expiration time.Duration) (ok bool) {
	if !d.connected {
		return
	}
	var err error
	if reflect.TypeOf(value).Kind() == reflect.Slice ||
		reflect.TypeOf(value).Kind() == reflect.Array ||
		reflect.TypeOf(value).Kind() == reflect.Map {
		var data string
		data, err = jsoniter.MarshalToString(value)
		if err == nil {
			err = d.client.Set(d.ctx, key, data, expiration).Err()
		}
	} else {
		err = d.client.Set(d.ctx, key, value, expiration).Err()
	}
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("[RedisCluster Error][Set Key=%s, Value=%+v] %s\n", key, value, err.Error())
		}
		ok = false
	} else {
		ok = true
	}
	return
}
func (d *RedisCluster) Delete(key string) (ok bool) {
	if !d.connected {
		return
	}
	effected, err := d.client.Del(d.ctx, key).Result()
	if effected > 0 {
		ok = true
	}
	if err != nil && err != redis.Nil {
		fmt.Printf("[RedisCluster Error][Delete Key=%s] %s\n", key, err.Error())
	}
	return
}
