package lache

import (
	"testing"
	"time"
)

func TestNewItem(t *testing.T) {
	onCreated := ItemEventHandler(func(k string, event ItemEvent, item Item) {
		t.Log("onEvent, k=", k, ",event=", event, ",item.value=", item.value)
	})
	onUpdated := ItemEventHandler(func(k string, event ItemEvent, item Item) {
		t.Log("onEvent, k=", k, ",event=", event, ",item.value=", item.value)
	})
	onDeleted := ItemEventHandler(func(k string, event ItemEvent, item Item) {
		t.Log("onEvent, k=", k, ",event=", event, ",item.value=", item.value)
	})
	SubItemEvent("k01", ItemCreate, &onCreated)
	SubItemEvent("k01", ItemUpdate, &onUpdated)
	SubItemEvent("k01", ItemDelete, &onDeleted)
	SubItemEvent("k02", ItemDelete, &onDeleted)

	Set("k02", "test", time.Second*2)

	NewItem("k01", ItemOptions{
		Mode:     Interval,
		Group:    DefaultGroup,
		Interval: time.Second,
		IntervalHandler: func() interface{} {
			return time.Now().Second()
		},
	})
	go func() {
		time.Sleep(time.Second * 3)
		UnsubItemEvent("k01", ItemUpdate, &onUpdated)
		time.Sleep(time.Second * 10)
		Delete("k01")
	}()
	time.Sleep(time.Second * 30)
}
