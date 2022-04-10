package lache

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	rwMutex    *sync.RWMutex
	key        string
	value      interface{}
	createdAt  int64       //创建时间
	updatedAt  int64       //更新时间
	updatedCnt int         //更新次数
	Opts       ItemOptions //客户端配置项
}

// Set 更新值
func (item *Item) Set(v interface{}) {
	defer item.rwMutex.Unlock()
	item.rwMutex.Lock()
	item.value = v
	item.updatedAt = time.Now().UnixNano()
	item.updatedCnt++
}

// update 更新值
func (item *Item) update() {
	defer item.rwMutex.Unlock()
	item.rwMutex.Lock()
	if item.Opts.IntervalHandler != nil {
		item.value = item.Opts.IntervalHandler()
		fmt.Println("item.value update = ", item.value)
		item.updatedAt = time.Now().UnixNano()
		item.updatedCnt++
	}
}

// Get 读取
func (item *Item) Get() interface{} {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	return item.value
}

// GetString 读字符串
func (item *Item) GetString() (s string) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	s, _ = item.value.(string)
	return
}

// GetInt 读Int
func (item *Item) GetInt() (i int) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	i, _ = item.value.(int)
	return
}

// GetInt64 读Int64
func (item *Item) GetInt64() (i64 int64) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	i64, _ = item.value.(int64)
	return
}

// GetStringMap 读 map[string]interface{}
func (item *Item) GetStringMap() (sm map[string]interface{}) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	sm, _ = item.value.(map[string]interface{})
	return
}

// GetStringMapString 读 map[string]string
func (item *Item) GetStringMapString() (sms map[string]string) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	sms, _ = item.value.(map[string]string)
	return
}

// GetStringSlice 读 []string
func (item *Item) GetStringSlice() (ss []string) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	ss, _ = item.value.([]string)
	return
}

// GetIntSlice 读 []int
func (item *Item) GetIntSlice() (si []int) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	si, _ = item.value.([]int)
	return
}

// GetInt64Slice 读 []int
func (item *Item) GetInt64Slice() (si64 []int64) {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	si64, _ = item.value.([]int64)
	return
}
