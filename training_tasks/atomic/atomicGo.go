package atomic

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	value atomic.Int64
	max   int64
}

func (cont *Counter) Increment() bool {
	for {
		cur := cont.value.Load()
		if cur == cont.max {
			return false
		}
		if cont.value.CompareAndSwap(cur, cur+1) {
			return true
		}
	}
}

func (cont *Counter) Value() int64 {
	return cont.value.Load()
}

func main() {
	var wg sync.WaitGroup
	var seccesses atomic.Int64
	counter := Counter{max: 1000}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				ok := counter.Increment()
				if ok {
					seccesses.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Результат счетчика: %d\nУспехи: %d\nНеудачи: %d", counter.Value(), seccesses.Load(), 10000-seccesses.Load())
}
