package lache

import (
	"sync"
	"time"
)

type (
	ItemMode int

	Item[K any, V any] struct {
		Mode       ItemMode      //模式 interval 循环更新模式, expire 过期模式
		Group      string        // 分组
		Interval   time.Duration // 循环更新时间
		Expires    time.Duration // 有效期限
		UpdatedAt  int64         //更新时间
		CreatedAt  int64         //创建时间
		UpdatedCnt int           //更新次数
		rwMutex    *sync.RWMutex
		key        K
		value      V
		updateFunc UpdateFunc
	}

	ItemOptions struct {
		Mode       ItemMode      //模式 interval 循环更新模式, expire 过期模式
		Group      string        // 分组
		Interval   time.Duration // 循环更新时间
		Expires    time.Duration // 有效期限
		UpdateFunc UpdateFunc
	}

	OptionFunc func(opt *ItemOptions)

	UpdateFunc func() any
)

const (
	DefaultGroup = "Default"

	Interval ItemMode = 0 // 定时更新模式
	Expire   ItemMode = 1 // 过期模式
	Manual   ItemMode = 2 //手动模式
)

// Set 更新值
func (item *Item[K, V]) Set(v V) {
	defer item.rwMutex.Unlock()
	item.rwMutex.Lock()
	item.value = v
	item.UpdatedAt = time.Now().UnixNano()
	item.UpdatedCnt++
}

// Update 更新值
func (item *Item[K, V]) Update() bool {
	defer item.rwMutex.Unlock()
	item.rwMutex.Lock()
	if item.updateFunc != nil {
		item.value, _ = item.updateFunc().(V)
		item.UpdatedAt = time.Now().UnixNano()
		item.UpdatedCnt++
		return true
	}
	return false
}

// Get 读取
func (item *Item[K, V]) Get() V {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	return item.value
}

// GetItem 读取
func (item *Item[K, V]) GetItem() Item[K, V] {
	defer item.rwMutex.RUnlock()
	item.rwMutex.RLock()
	return *item
}

func WithMode(m ItemMode) OptionFunc {
	return func(opt *ItemOptions) {
		opt.Mode = m
	}
}
func WithInterval(d time.Duration) OptionFunc {
	return func(opt *ItemOptions) {
		opt.Interval = d
	}
}
func WithExpires(d time.Duration) OptionFunc {
	return func(opt *ItemOptions) {
		opt.Expires = d
	}
}
func WithGroup(g string) OptionFunc {
	return func(opt *ItemOptions) {
		opt.Group = g
	}
}
func WithUpdateFunc(f UpdateFunc) OptionFunc {
	return func(opt *ItemOptions) {
		opt.UpdateFunc = f
	}
}
