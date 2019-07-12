// Package safemap provides a safe map of string -> interface
package safemap

import deadlock "github.com/sasha-s/go-deadlock"

// SafeMap provides a map safe for concurrent updates
type SafeMap struct {
	deadlock.RWMutex
	data map[string]interface{}
}

type SafeMapItem struct {
	Key   string
	Value interface{}
}

func New() *SafeMap {
	return &SafeMap{
		data: make(map[string]interface{}),
	}
}

func (sm *SafeMap) Keys() []string {
	sm.RLock()
	defer sm.RUnlock()

	var keys []string
	for k := range sm.data {
		keys = append(keys, k)
	}
	return keys
}

func (sm *SafeMap) Iter() <-chan SafeMapItem {
	c := make(chan SafeMapItem, 10)

	f := func() {
		sm.RLock()
		defer sm.RUnlock()
		for k, v := range sm.data {
			c <- SafeMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}

func (sm *SafeMap) Insert(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.data[key] = value
}

func (sm *SafeMap) Delete(key string) {
	sm.Lock()
	defer sm.Unlock()
	delete(sm.data, key)
}

func (sm *SafeMap) Find(key string) (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	v, ok := sm.data[key]
	return v, ok
}

func (sm *SafeMap) Len() int {
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.data)
}

func (sm *SafeMap) Update(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.data[key] = value
}
