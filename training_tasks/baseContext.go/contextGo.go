package basecontextgogo

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: stopped: %v\n", id, ctx.Err())
			return
		case <-time.After(200 * time.Millisecond):
			fmt.Printf("worker %d: tick\n", id)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go Worker(ctx, i, &wg)
	}
	time.Sleep(1 * time.Second)
	cancel()

	wg.Wait()
	fmt.Println("all workers stopped")
}
