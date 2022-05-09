package lache

import (
	"reflect"
	"sync"
	"time"
)

var items map[any]any
var itemsRWMutex *sync.RWMutex

func init() {
	//初始化读写锁
	itemsRWMutex = new(sync.RWMutex)
	//初始化map队列
	items = make(map[any]any)
	//计时器
	go tick()
}

func tick() {
	tickInterval := time.Second
	timer := time.NewTimer(tickInterval)
	var updatedAt int64
	var interval, expires time.Duration
	var mode ItemMode
	for {
		select {
		case <-timer.C:
			itemsRWMutex.Lock() //写锁
			nano := time.Now().UnixNano()
			for key, item := range items {
				mode, _ = reflect.ValueOf(item).Elem().FieldByName("Mode").Interface().(ItemMode)
				updatedAt, _ = reflect.ValueOf(item).Elem().FieldByName("UpdatedAt").Interface().(int64)
				if mode == Expire { // 过期模式
					expires, _ = reflect.ValueOf(item).Elem().FieldByName("Expires").Interface().(time.Duration)
					if updatedAt+int64(expires) < nano {
						pubItemEvent(key, ItemDelete, reflect.ValueOf(item).Elem().Interface())
						ClearItemEvent(key)
						delete(items, key)
					}
				} else if mode == Interval { // 定时更新模式
					interval, _ = reflect.ValueOf(item).Elem().FieldByName("Interval").Interface().(time.Duration)
					if updatedAt+int64(interval) < nano {
						callRes := reflect.ValueOf(item).MethodByName("Update").Call(nil)
						if len(callRes) > 0 {
							if updateRes, ok := callRes[0].Interface().(bool); ok && updateRes {
								pubItemEvent(key, ItemUpdate, reflect.ValueOf(item).Elem().Interface())
							}
						}
					}
				}
			}
			itemsRWMutex.Unlock() //释放写锁
			timer.Reset(tickInterval)
		}
	}
}

// Set 设置kv
func Set(k string, v string, expires time.Duration) bool {
	return SetKV[string, string](k, v, expires)
}

// Get 获取v
func Get(k string) string {
	return GetKV[string, string](k)
}

// Update 更新v
func Update(k, v string) bool {
	return UpdateKV[string, string](k, v)
}

// SetKV 设置kv
func SetKV[K any, V any](k K, v V, expires time.Duration) bool {
	_, ok := NewItemKV[K, V](k, v, WithExpires(expires))
	return ok
}

// GetKV 获取
func GetKV[K any, V any](k K) (v V) {
	defer itemsRWMutex.RUnlock() //释放读锁
	itemsRWMutex.RLock()         //读锁
	item, exist := items[k]
	if !exist {
		return
	}
	values := reflect.ValueOf(item).MethodByName("Get").Call(nil)
	if len(values) > 0 {
		v, _ = values[0].Interface().(V)
	}
	return
}

// UpdateKV 更新
func UpdateKV[K any, V any](k K, v V) bool {
	defer itemsRWMutex.RUnlock() //释放读锁
	itemsRWMutex.RLock()         //读锁
	item, exist := items[k]
	if !exist {
		return false
	}
	reflect.ValueOf(item).MethodByName("Set").Call([]reflect.Value{reflect.ValueOf(v)})
	pubItemEvent(k, ItemUpdate, reflect.ValueOf(item).Elem().Interface())
	return true
}

// Delete 删除
func Delete(k any) {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	if v, ok := items[k]; ok {
		pubItemEvent(k, ItemDelete, reflect.ValueOf(v).Elem().Interface())
		delete(items, k)
		ClearItemEvent(k)
	}
}

// NewItemKV 新的缓存项
func NewItemKV[K any, V any](k K, v V, opts ...OptionFunc) (item *Item[K, V], ok bool) {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	if _, exist := items[k]; exist {
		return
	}
	nano := time.Now().UnixNano()

	defaultOpts := ItemOptions{
		Mode:     Expire,
		Group:    DefaultGroup,
		Interval: time.Hour,
		Expires:  time.Hour,
	}
	for _, opt := range opts {
		opt(&defaultOpts)
	}
	item = &Item[K, V]{
		key:        K(k),
		value:      V(v),
		Group:      defaultOpts.Group,
		Mode:       defaultOpts.Mode,
		Expires:    defaultOpts.Expires,
		Interval:   defaultOpts.Interval,
		updateFunc: defaultOpts.UpdateFunc,
		CreatedAt:  nano,
		UpdatedAt:  nano,
		rwMutex:    new(sync.RWMutex),
	}
	items[k] = item
	if item.Mode == Interval || item.Mode == Manual {
		item.Update()
	}
	pubItemEvent(k, ItemCreate, *item)
	return
}
