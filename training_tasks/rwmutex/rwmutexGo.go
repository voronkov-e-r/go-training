package rwmutex

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Entry struct {
	value     string
	expiresAt time.Time
}

type Cash struct {
	data map[string]Entry
	mu   sync.RWMutex
}

func (cash *Cash) Get(key string) (val string, ok bool) {
	cash.mu.RLock()
	obj, err := cash.data[key]
	cash.mu.RUnlock()
	val = obj.value
	ok = err
	now := time.Now()
	if obj.expiresAt.Before(now) {
		cash.mu.Lock()
		defer cash.mu.Unlock()
		if cur, status := cash.data[key]; status && cur.expiresAt.Before(time.Now()) {
			delete(cash.data, key)
			ok = false
		}
		val = ""
	}
	return
}

func (cash *Cash) Set(key, value string, ttl time.Duration) {
	cash.mu.Lock()
	defer cash.mu.Unlock()
	expire := time.Now().Add(ttl)
	cash.data[key] = Entry{value, expire}
}

func (cash *Cash) Delete(key string) {
	cash.mu.Lock()
	defer cash.mu.Unlock()
	delete(cash.data, key)
}

func (cash *Cash) Len() int {
	len := 0
	now := time.Now()
	cash.mu.RLock()
	keys := []string{}
	for key, obj := range cash.data {
		if obj.expiresAt.Before(now) {
			keys = append(keys, key)
		} else {
			len += 1
		}
	}
	cash.mu.RUnlock()

	cash.mu.Lock()
	defer cash.mu.Unlock()
	for _, key := range keys {
		if cash.data[key].expiresAt.Before(time.Now()) {
			delete(cash.data, key)
		}
	}

	return len
}

func main() {
	var wg sync.WaitGroup
	var misses atomic.Int64
	var cash Cash
	cash.data = map[string]Entry{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cash.Set(fmt.Sprintf("k%d_%d", i, j), "", time.Duration(500*time.Millisecond))
			}
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				firstIndex := rand.Intn(10)
				secondIndex := rand.Intn(100)
				_, ok := cash.Get(fmt.Sprintf("k%d_%d", firstIndex, secondIndex))
				if !ok {
					misses.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	startLen := cash.Len()
	time.Sleep(600 * time.Millisecond)
	endLen := cash.Len()
	heats := 2000 - misses.Load()
	fmt.Printf("Hits: %d\nMisses: %d\nStart Len: %d\nEnd Len: %d", heats, misses.Load(), startLen, endLen)
}
