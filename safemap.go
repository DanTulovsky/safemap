// Package safemap provides a safe map of string -> interface
package safemap

import deadlock "github.com/sasha-s/go-deadlock"

type safeMap struct {
	deadlock.RWMutex
	data map[string]interface{}
}

type safeMapItem struct {
	Key   string
	Value interface{}
}

func NewSafeMap() *safeMap {
	return &safeMap{
		data: make(map[string]interface{}),
	}
}

func (sm *safeMap) Keys() []string {
	sm.RLock()
	defer sm.RUnlock()

	var keys []string
	for k := range sm.data {
		keys = append(keys, k)
	}
	return keys
}

func (sm *safeMap) Iter() <-chan safeMapItem {
	c := make(chan safeMapItem, 10)

	f := func() {
		sm.RLock()
		defer sm.RUnlock()
		for k, v := range sm.data {
			c <- safeMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}

func (sm *safeMap) Insert(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.data[key] = value
}

func (sm *safeMap) Delete(key string) {
	sm.Lock()
	defer sm.Unlock()
	delete(sm.data, key)
}

func (sm *safeMap) Find(key string) (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	v, ok := sm.data[key]
	return v, ok
}

func (sm *safeMap) Len() int {
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.data)
}

func (sm *safeMap) Update(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.data[key] = value
}
