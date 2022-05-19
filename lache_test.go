package lache

import (
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhouhp1295/lache/driver"
	"testing"
	"time"
)

type TestData struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

func (t *TestData) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(*t)
}

func (t *TestData) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, t)
}

func TestLocal(t *testing.T) {
	/*
		client := New(Local, driver.LocalOptions{})
		client.Set("key01", 123, driver.NotExpired)
		t.Log(client.Get("key01"))
		//基础类型测试
		var result01 int
		t.Log("key01", client.GetT("key01", &result01), result01)
		//指针测试
		client.Set("key02", TestData{Id: 2, Value: "2"}, driver.NotExpired)
		var result02 TestData
		t.Log("key02", client.GetT("key02", &result02), result02)
		//new指针测试
		client.Set("key03", &TestData{Id: 3, Value: "3"}, driver.NotExpired)
		result03 := new(TestData)
		t.Log("key03", client.GetT("key03", result03), result03)
		//数组测试
		client.Set("key04", []int{1, 2, 3}, driver.NotExpired)
		var result04 []int
		t.Log("key04", client.GetT("key04", &result04), result04)
		//map测试
		client.Set("key05", map[string]string{"a": "a", "b": "b"}, driver.NotExpired)
		var result05 map[string]string
		t.Log("key05", client.GetT("key05", &result05), result05)

		client.Set("key06", 12345, time.Second*5)
		for i := 0; i < 10; i++ {
			t.Log(client.Get("key06"))
			time.Sleep(time.Second)
		}
	*/
}

func TestRedis(t *testing.T) {
	client := New(Redis, redis.Options{
		Addr: "127.0.0.1:6379",
	})
	//基础类型测试
	client.Set("key01", 123.92, driver.NotExpired)
	t.Log(client.Get("key01"))
	var result01 float32
	t.Log("key01", client.GetT("key01", &result01), result01)
	//指针测试
	client.Set("key02", &TestData{Id: 2, Value: "2"}, driver.NotExpired)
	var result02 TestData
	t.Log("key02", client.GetT("key02", &result02), result02)
	//new指针测试
	client.Set("key03", &TestData{Id: 3, Value: "3"}, driver.NotExpired)
	result03 := new(TestData)
	t.Log("key03", client.GetT("key03", result03), result03)
	//数组测试
	client.Set("key04", []int{1, 2, 3}, driver.NotExpired)
	var result04 []int
	t.Log("key04", client.GetT("key04", &result04), result04)
	client.Set("key04_1", [3]int{1, 2, 3}, driver.NotExpired)
	var result04_1 [3]int
	t.Log("key04_1", client.GetT("key04", &result04_1), result04_1)
	//map测试
	client.Set("key05", map[string]string{"a": "a", "b": "b"}, driver.NotExpired)
	var result05 map[string]string
	t.Log("key05", client.GetT("key05", &result05), result05)

	client.Set("key06", 12345, time.Second*5)
	for i := 0; i < 10; i++ {
		t.Log(client.Get("key06"))
		time.Sleep(time.Second)
	}
}
