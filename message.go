package lache

import (
	"sync"
)

type SubHandler func(params ...interface{})
type ItemEventHandler func(k string, event ItemEvent, item Item)

type ItemEvent string

const ItemCreate ItemEvent = "CREATE"
const ItemUpdate ItemEvent = "UPDATE"
const ItemDelete ItemEvent = "DELETE"

var handlersRWMutex *sync.RWMutex
var handlers map[string][]*SubHandler
var itemEventHandlersRWMutex *sync.RWMutex
var itemEventHandlers map[string]map[ItemEvent][]*ItemEventHandler

func init() {
	handlersRWMutex = new(sync.RWMutex)
	handlers = make(map[string][]*SubHandler)
	itemEventHandlersRWMutex = new(sync.RWMutex)
	itemEventHandlers = make(map[string]map[ItemEvent][]*ItemEventHandler)
}

// Sub 订阅
func Sub(topic string, handler *SubHandler) {
	defer handlersRWMutex.Unlock() //释放写锁
	handlersRWMutex.Lock()         //写锁
	if _, exist := handlers[topic]; !exist {
		handlers[topic] = make([]*SubHandler, 0)
	}
	handlers[topic] = append(handlers[topic], handler)
}

// Unsub 取消订阅
func Unsub(topic string, handler *SubHandler) {
	defer handlersRWMutex.Unlock() //释放写锁
	handlersRWMutex.Lock()         //写锁
	if _, exist := handlers[topic]; exist {
		for i, _ := range handlers[topic] {
			handlers[topic] = append(handlers[topic][:i], handlers[topic][i+1:]...)
			return
		}
	}
}

// Pub 发布
func Pub(topic string, params ...interface{}) {
	defer handlersRWMutex.RUnlock() //释放读锁
	handlersRWMutex.RLock()         //读锁
	if _, exist := handlers[topic]; exist {
		for _, f := range handlers[topic] {
			(*f)(params)
		}
	}
}

// pubItemEvent 发布缓存项事件
func pubItemEvent(k string, event ItemEvent, item Item) {
	defer itemEventHandlersRWMutex.RUnlock() //释放读锁
	itemEventHandlersRWMutex.RLock()         //读锁
	if _, exist := itemEventHandlers[k]; exist {
		if data, exist2 := itemEventHandlers[k][event]; exist2 {
			for _, f := range data {
				(*f)(k, event, item)
				return
			}
		}
	}
}

// SubItemEvent 订阅缓存项事件
func SubItemEvent(k string, event ItemEvent, handler *ItemEventHandler) {
	defer itemEventHandlersRWMutex.Unlock() //释放写锁
	itemEventHandlersRWMutex.Lock()         //写锁
	if _, exist := itemEventHandlers[k]; !exist {
		itemEventHandlers[k] = make(map[ItemEvent][]*ItemEventHandler)
	}
	if _, exist := itemEventHandlers[k][event]; !exist {
		itemEventHandlers[k][event] = make([]*ItemEventHandler, 0)
	}
	itemEventHandlers[k][event] = append(itemEventHandlers[k][event], handler)
}

// UnsubItemEvent 取消订阅缓存项事件
func UnsubItemEvent(k string, event ItemEvent, handler *ItemEventHandler) {
	defer itemEventHandlersRWMutex.Unlock() //释放写锁
	itemEventHandlersRWMutex.Lock()         //写锁
	if _, exist := itemEventHandlers[k]; exist {
		if _, exist2 := itemEventHandlers[k][event]; exist2 {
			for i, _ := range itemEventHandlers[k][event] {
				itemEventHandlers[k][event] = append(itemEventHandlers[k][event][:i], itemEventHandlers[k][event][i+1:]...)
				return
			}
		}
	}
}
