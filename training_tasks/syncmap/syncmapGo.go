package syncmap

import (
	"sync"
)

type MutexMap struct {
	data map[string]int
	mu   sync.Mutex
}

func (mp *MutexMap) Get(key string) (int, bool) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	val, ok := mp.data[key]
	return val, ok
}

func (mp *MutexMap) Set(key string, value int) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.data[key] = value
}

type RWMutexMap struct {
	data map[string]int
	mu   sync.RWMutex
}

func (mp *RWMutexMap) Get(key string) (int, bool) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	val, ok := mp.data[key]
	return val, ok
}

func (mp *RWMutexMap) Set(key string, value int) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.data[key] = value
}

type SyncMap struct {
	data sync.Map
}

func (mp *SyncMap) Get(key any) (any, bool) {
	val, ok := mp.data.Load(key)
	return val, ok
}

func (mp *SyncMap) Set(key, value any) {
	mp.data.Store(key, value)
}

func main() {

}
