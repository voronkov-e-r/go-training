package pool

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func main() {
	var wg sync.WaitGroup

	handleRequest := func(id int) string {
		buf := bufPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufPool.Put(buf)
		}()
		fmt.Fprintf(buf, "request #%d processed\n", id)
		time.Sleep(1 * time.Millisecond)
		return strings.Clone(buf.String())
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleRequest(i)
		}()
	}
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Время выполнения: %d", duration.Milliseconds())
}
