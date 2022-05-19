package event

import "sync"

type Hub struct {
	rwMutex  *sync.RWMutex
	handlers map[string][]*Handler
}

type Handler func(topic string, params ...any)

func New() *Hub {
	hub := new(Hub)
	hub.rwMutex = new(sync.RWMutex)
	hub.handlers = make(map[string][]*Handler)
	return hub
}

// Sub 订阅
func (hub *Hub) Sub(topic string, handler *Handler) {
	defer hub.rwMutex.Unlock() //释放写锁
	hub.rwMutex.Lock()         //写锁
	if _, exist := hub.handlers[topic]; !exist {
		hub.handlers[topic] = make([]*Handler, 0)
	}
	hub.handlers[topic] = append(hub.handlers[topic], handler)
}

// Unsub 取消订阅
func (hub *Hub) Unsub(topic string, handler *Handler) {
	defer hub.rwMutex.Unlock() //释放写锁
	hub.rwMutex.Lock()         //写锁
	if _, exist := hub.handlers[topic]; exist {
		for i, h := range hub.handlers[topic] {
			if h == handler {
				hub.handlers[topic] = append(hub.handlers[topic][:i], hub.handlers[topic][i+1:]...)
				return
			}
		}
	}
}

// Pub 发布
func (hub *Hub) Pub(topic string, params ...any) {
	defer hub.rwMutex.RUnlock() //释放读锁
	hub.rwMutex.RLock()         //读锁
	if _, exist := hub.handlers[topic]; exist {
		for _, f := range hub.handlers[topic] {
			(*f)(topic, params...)
		}
	}
}
