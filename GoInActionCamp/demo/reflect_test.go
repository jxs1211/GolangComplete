package reflect_demo

import (
	"reflect"
	"sync"
	"testing"
)

type any interface{}

func TestReflectTypeSlice(t *testing.T) {
	s := []int{1, 2, 3}
	val := reflect.ValueOf(s)
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		res = append(res, val.Index(i).Interface())
	}
	t.Log(res)
}

type SafeMap struct {
	m  map[string]string
	mu sync.RWMutex
}

func (m *SafeMap) Get(key string) string {
	// 只读需要枷锁，因为有其他goroutine在同时写
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[key]
}

func (m *SafeMap) Set(key, val string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = val
}
