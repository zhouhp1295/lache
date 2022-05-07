package lache

import "sync"

type (
	ItemEvent        string
	SubHandler       func(params ...any)
	ItemEventHandler func(k any, event ItemEvent, item interface{})
)

const (
	ItemCreate ItemEvent = "CREATE"
	ItemUpdate ItemEvent = "UPDATE"
	ItemDelete ItemEvent = "DELETE"
)

var (
	handlersRWMutex          *sync.RWMutex
	handlers                 map[string][]*SubHandler
	itemEventHandlersRWMutex *sync.RWMutex
	itemEventHandlers        map[any]map[ItemEvent][]*ItemEventHandler
)

func init() {
	handlersRWMutex = new(sync.RWMutex)
	handlers = make(map[string][]*SubHandler)
	itemEventHandlersRWMutex = new(sync.RWMutex)
	itemEventHandlers = make(map[any]map[ItemEvent][]*ItemEventHandler)
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
func Pub(topic string, params ...any) {
	defer handlersRWMutex.RUnlock() //释放读锁
	handlersRWMutex.RLock()         //读锁
	if _, exist := handlers[topic]; exist {
		for _, f := range handlers[topic] {
			(*f)(params)
		}
	}
}

// pubItemEvent 发布缓存项事件
func pubItemEvent(k any, event ItemEvent, item interface{}) {
	defer itemEventHandlersRWMutex.RUnlock() //释放读锁
	itemEventHandlersRWMutex.RLock()         //读锁
	if _, exist := itemEventHandlers[k]; exist {
		if data, exist2 := itemEventHandlers[k][event]; exist2 {
			for _, f := range data {
				(*f)(k, event, item)
			}
		}
	}
}

// SubItemEvent 订阅缓存项事件
func SubItemEvent(k any, event ItemEvent, handler *ItemEventHandler) {
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
func UnsubItemEvent(k any, event ItemEvent, handler *ItemEventHandler) {
	defer itemEventHandlersRWMutex.Unlock() //释放写锁
	itemEventHandlersRWMutex.Lock()         //写锁
	if _, exist := itemEventHandlers[k]; exist {
		if _, exist2 := itemEventHandlers[k][event]; exist2 {
			for i, h := range itemEventHandlers[k][event] {
				if h == handler {
					itemEventHandlers[k][event] = append(itemEventHandlers[k][event][:i], itemEventHandlers[k][event][i+1:]...)
					return
				}
			}
		}
	}
}

// ClearItemEvent 清空
func ClearItemEvent(k any) {
	defer itemEventHandlersRWMutex.Unlock() //释放写锁
	itemEventHandlersRWMutex.Lock()         //写锁
	if _, exist := itemEventHandlers[k]; exist {
		delete(itemEventHandlers, k)
	}
}
