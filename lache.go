package lache

import (
	"sync"
	"time"
)

var items map[string]*Item
var itemsRWMutex *sync.RWMutex
var defaultIntervalItemOpts ItemOptions //默认的创建配置

func init() {
	//初始化读写锁
	itemsRWMutex = new(sync.RWMutex)
	//初始化默认选项
	defaultIntervalItemOpts = ItemOptions{
		Group:   DefaultGroup,
		Mode:    Expire,
		Expires: time.Hour,
	}
	//初始化map队列
	items = make(map[string]*Item)
	//开启timer
	go tick()
}

func tick() {
	interval := time.Second
	timer := time.NewTimer(interval)
	for {
		select {
		case <-timer.C:
			itemsRWMutex.Lock() //写锁
			nano := time.Now().UnixNano()
			for key, item := range items {
				if item.Opts.Mode == Interval {
					// 自动更新模式
					if item.updatedAt+int64(item.Opts.Interval) < nano {
						item.update()
						pubItemEvent(key, ItemUpdate, *item)
					}
				} else if item.Opts.Mode == Expire {
					// 过期模式
					if item.updatedAt+int64(item.Opts.Expires) < nano {
						pubItemEvent(key, ItemDelete, *item)
						delete(items, key)
					}
				}
			}
			itemsRWMutex.Unlock() //释放写锁
			timer.Reset(interval)
		}
	}
}

// NewDefaultItem 新建缓存项
func NewDefaultItem(k string) (*Item, bool) {
	return NewItem(k, defaultIntervalItemOpts)
}

// NewItem 新建缓存项
func NewItem(k string, itemOpts ItemOptions) (*Item, bool) {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	if _, exist := items[k]; exist {
		return nil, false
	}
	nano := time.Now().UnixNano()
	item := &Item{
		Opts:      itemOpts,
		key:       k,
		createdAt: nano,
		updatedAt: nano,
		rwMutex:   new(sync.RWMutex),
	}
	if itemOpts.Mode == Interval && itemOpts.IntervalHandler != nil {
		item.value = itemOpts.IntervalHandler()
	}
	items[k] = item
	pubItemEvent(k, ItemCreate, *item)
	return item, true
}

// NewItemWithOption 新建缓存项
func NewItemWithOption(k string, opts ...Option) (*Item, bool) {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	if _, exist := items[k]; exist {
		return nil, false
	}
	itemOpts := defaultIntervalItemOpts
	for _, opt := range opts {
		opt(&itemOpts)
	}
	nano := time.Now().UnixNano()
	item := &Item{
		Opts:      itemOpts,
		key:       k,
		createdAt: nano,
		updatedAt: nano,
		rwMutex:   new(sync.RWMutex),
	}
	if itemOpts.Mode == Interval && itemOpts.IntervalHandler != nil {
		item.value = itemOpts.IntervalHandler()
	}
	items[k] = item
	pubItemEvent(k, ItemCreate, *item)
	return item, true
}

// Set 设置kv
func Set(k, v string, expires time.Duration) bool {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	if _, exist := items[k]; exist {
		return false
	}
	nano := time.Now().UnixNano()
	item := &Item{
		rwMutex: new(sync.RWMutex),
		Opts: ItemOptions{
			Mode:    Expire,
			Group:   DefaultGroup,
			Expires: expires,
		},
		createdAt: nano,
		updatedAt: nano,
		key:       k,
		value:     v,
	}
	items[k] = item
	pubItemEvent(k, ItemCreate, *item)
	return true
}

// Get 获取
func Get(k string) string {
	defer itemsRWMutex.RUnlock() //释放读锁
	itemsRWMutex.RLock()         //读锁
	item, exist := items[k]
	if !exist {
		return ""
	}
	return item.GetString()
}

// Update 更新
func Update(k, v string) bool {
	defer itemsRWMutex.RUnlock() //释放读锁
	itemsRWMutex.RLock()         //读锁
	item, exist := items[k]
	if !exist {
		return false
	}
	item.Set(v)
	pubItemEvent(k, ItemUpdate, *item)
	return true
}

// Delete 删除
func Delete(k string) {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	if v, ok := items[k]; ok {
		pubItemEvent(k, ItemDelete, *v)
		delete(items, k)
	}

}

// Clear 清空
func Clear() {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	for k, v := range items {
		pubItemEvent(k, ItemDelete, *v)
	}
	items = make(map[string]*Item)
}

// DeleteGroup 根据分组删除
func DeleteGroup(group string) {
	defer itemsRWMutex.Unlock() //释放写锁
	itemsRWMutex.Lock()         //写锁
	for k, v := range items {
		if v.Opts.Group == group {
			pubItemEvent(k, ItemDelete, *v)
			delete(items, k)
		}
	}
}
