package cond

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Buffer struct {
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	data     []int
	capacity int
}

func NewBuffer(cap int) *Buffer {
	b := &Buffer{capacity: cap}
	b.notEmpty = sync.NewCond(&b.mu)
	b.notFull = sync.NewCond(&b.mu)
	return b
}

func (ch *Buffer) Put(val int) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	for len(ch.data) == ch.capacity {
		ch.notFull.Wait()
	}
	ch.data = append(ch.data, val)
	ch.notEmpty.Signal()
}

func (ch *Buffer) Get() int {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	for len(ch.data) == 0 {
		ch.notEmpty.Wait()
	}
	v := ch.data[0]
	ch.data = ch.data[1:]
	ch.notFull.Signal()
	return v
}

func main() {
	ch := NewBuffer(10)
	var wg sync.WaitGroup
	var totalSum atomic.Int64

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				ch.Put(base*100 + j)
			}
		}(i)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				val := ch.Get()
				totalSum.Add(int64(val))
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Итоговая сумма: %d\n", totalSum.Load())
}
