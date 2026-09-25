package once

import (
	"fmt"
	"sync"
	"time"
)

type DB struct {
	conn string
	once sync.Once
}

func (db *DB) Connect() string {
	db.once.Do(func() {
		time.Sleep(100 * time.Millisecond)
		db.conn = "connected-to-db"
	})
	return db.conn
}

func main() {
	var db DB
	var wg sync.WaitGroup
	var mu sync.Mutex
	unique := make(map[string]int)

	start := time.Now()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := db.Connect()
			mu.Lock()
			defer mu.Unlock()
			_, ok := unique[res]
			if !ok {
				unique[res] = 1
			} else {
				unique[res] += 1
			}
		}()
	}

	wg.Wait()

	duration := time.Since(start).Milliseconds()
	fmt.Printf("Время выполнения: %dms\nУникальные значения: %d", duration, len(unique))
}
