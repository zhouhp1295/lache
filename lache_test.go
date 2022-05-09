package lache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewItem(t *testing.T) {
	onEvent := ItemEventHandler(func(k any, event ItemEvent, item interface{}) {
		if item0, ok := item.(Item[string, string]); ok {
			t.Logf("onEvent, k= %v, event = %v, item = %v", k, event, item0.Get())
		}
	})
	onEvent2 := ItemEventHandler(func(k any, event ItemEvent, item interface{}) {
		if item0, ok := item.(Item[string, string]); ok {
			t.Logf("onEvent2, k= %v, event = %v, item = %v", k, event, item0.Get())
		}
	})
	SubItemEvent("k01", ItemCreate, &onEvent)
	SubItemEvent("k01", ItemDelete, &onEvent)
	SubItemEvent("k01", ItemUpdate, &onEvent)
	SubItemEvent("k02", ItemUpdate, &onEvent)
	SubItemEvent("k02", ItemUpdate, &onEvent2)
	SubItemEvent("k02", ItemDelete, &onEvent)
	SubItemEvent("k04", ItemUpdate, &onEvent)
	SubItemEvent("k04", ItemDelete, &onEvent)
	Set("k01", "test01", time.Second*2)

	NewItemKV[string, string]("k02", "0", WithMode(Interval), WithInterval(time.Second), WithUpdateFunc(func() any {
		return fmt.Sprintf("%d", time.Now().Second())
	}))
	//NewItemKV[string, string]("k02", "0", WithMode(Interval), WithInterval(time.Second))

	Set("k04", "test04", time.Hour)
	t.Log(Get("k04"))
	Update("k04", "test04-update")

	go func() {
		time.Sleep(time.Second * 2)
		UnsubItemEvent("k02", ItemUpdate, &onEvent)
		Update("k04", "test04-update-update")
		time.Sleep(time.Second * 5)
		Delete("k02")
		time.Sleep(time.Second * 2)
	}()

	time.Sleep(time.Second * 15)
}
