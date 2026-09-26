package hardcontext

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func fetch(ctx context.Context, url string) (string, error) {
	delay := time.Duration(100+rand.Intn(250)) * time.Millisecond
	select {
	case <-time.After(delay):
		return "data from " + url, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func fetchAll(ctx context.Context, urls []string) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		data string
		err  error
	}
	ch := make(chan result, len(urls))

	for _, url := range urls {
		go func(u string) {
			data, err := fetch(ctx, u)
			select {
			case ch <- result{data, err}:
			case <-ctx.Done():

			}
		}(url)
	}

	var lastErr error
	for i := 0; i < len(urls); i++ {
		select {
		case r := <-ch:
			if r.err == nil {
				return []string{r.data}, nil
			}
			lastErr = r.err
		case <-ctx.Done():
			if lastErr != nil {
				return nil, lastErr
			}
			return nil, ctx.Err()
		}
	}
	return nil, lastErr
}

func main() {
	before := runtime.NumGoroutine()
	fmt.Println("Тест 1:")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	urls := []string{"192.34.56.01", "165.78.10.23", "144.65.7.21", "122.45.67.5"}

	res, err := fetchAll(ctx, urls)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Результаты:")
		for _, str := range res {
			fmt.Println(str)
		}
	}

	fmt.Println("Тест 2:")
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(80 * time.Millisecond)
		cancel()
	}()

	res, err = fetchAll(ctx, urls)

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Результаты:")
		for _, str := range res {
			fmt.Println(str)
		}
	}

	fmt.Println("Тест 3:")
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	res, err = fetchAll(ctx, urls)

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Результаты:")
		for _, str := range res {
			fmt.Println(str)
		}
	}
	time.Sleep(500 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Printf("Goroutines: before=%d, after=%d\n", before, after)
}
