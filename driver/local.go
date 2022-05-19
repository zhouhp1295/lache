package driver

import (
	"reflect"
	"sync"
	"time"
)

const NotExpired = 0

// LocalOptions 初始化选项 [预留]
type LocalOptions struct {
	MaxSize int64
}

type localItem struct {
	value      any
	expiration int64
}

type Local struct {
	items   map[string]localItem
	rwMutex *sync.RWMutex
	Options *LocalOptions
}

func NewLocalDriver(options LocalOptions) *Local {
	driver := new(Local)
	driver.Options = &options
	driver.rwMutex = new(sync.RWMutex)
	driver.items = make(map[string]localItem)
	go driver.tick()
	return driver
}

func (d *Local) Get(key string) (result any, ok bool) {
	defer d.rwMutex.RUnlock() //释放读锁
	d.rwMutex.RLock()         //读锁
	item, exist := d.items[key]
	if !exist {
		return
	}
	if item.expiration == NotExpired || time.Now().UnixNano() < item.expiration {
		result = item.value
		ok = true
	}
	return
}

func (d *Local) GetT(key string, result any) (ok bool) {
	res, found := d.Get(key)
	if !found {
		return
	}
	if reflect.TypeOf(result).Kind() == reflect.Ptr {
		if reflect.TypeOf(res).Kind() == reflect.Ptr {
			if reflect.TypeOf(result).Elem().Kind() == reflect.TypeOf(res).Elem().Kind() {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(res).Elem())
				ok = true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.TypeOf(res).Kind() {
			reflect.ValueOf(result).Elem().Set(reflect.ValueOf(res))
			ok = true
		}
	}
	return
}

func (d *Local) Set(key string, value any, expiration time.Duration) (ok bool) {
	defer d.rwMutex.Unlock() //释放写锁
	d.rwMutex.Lock()         //写锁

	item := localItem{value: value, expiration: NotExpired}
	if expiration != NotExpired {
		item.expiration = time.Now().Add(expiration).UnixNano()
	}

	d.items[key] = item
	ok = true
	return
}

func (d *Local) Delete(key string) (ok bool) {
	defer d.rwMutex.Unlock() //释放写锁
	d.rwMutex.Lock()         //写锁
	ok = true
	delete(d.items, key)
	return
}

func (d *Local) tick() {
	tickInterval := time.Millisecond * 1000 // 1秒
	timer := time.NewTimer(tickInterval)
	for {
		select {
		case <-timer.C:
			nano := time.Now().UnixNano()
			d.rwMutex.RLock() //读锁
			for key, item := range d.items {
				//过期删除
				if nano > item.expiration {
					d.rwMutex.RUnlock() //释放读锁
					d.rwMutex.Lock()    //加写锁
					delete(d.items, key)
					d.rwMutex.Unlock() //释放写锁
					d.rwMutex.RLock()  //重新加读锁
				}
			}
			d.rwMutex.RUnlock() //释放读锁
			timer.Reset(tickInterval)
		}
	}
}
