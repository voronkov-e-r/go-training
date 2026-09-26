package syncmap

import (
	"strconv"
	"testing"
)

func BenchmarkMutexReadHeavy(b *testing.B) {
	c := &MutexMap{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				c.Set("k"+strconv.Itoa(i%100), i)
			} else {
				c.Get("k" + strconv.Itoa(i%100))
			}
			i++
		}
	})
}

func BenchmarkRWMutexReadHeavy(b *testing.B) {
	c := &RWMutexMap{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				c.Set("k"+strconv.Itoa(i%100), i)
			} else {
				c.Get("k" + strconv.Itoa(i%100))
			}
			i++
		}
	})
}

func BenchmarkSyncReadHeavy(b *testing.B) {
	c := &SyncMap{}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				c.Set("k"+strconv.Itoa(i%100), i)
			} else {
				c.Get("k" + strconv.Itoa(i%100))
			}
			i++
		}
	})
}

func BenchmarkMutexWriteHeavy(b *testing.B) {
	c := &MutexMap{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				c.Get("k" + strconv.Itoa(i%100))
			} else {
				c.Set("k"+strconv.Itoa(i%100), i)
			}
			i++
		}
	})
}

func BenchmarkRWMutexWriteHeavy(b *testing.B) {
	c := &RWMutexMap{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				c.Get("k" + strconv.Itoa(i%100))
			} else {
				c.Set("k"+strconv.Itoa(i%100), i)
			}
			i++
		}
	})
}

func BenchmarkSyncWriteHeavy(b *testing.B) {
	c := &SyncMap{}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				c.Get("k" + strconv.Itoa(i%100))
			} else {
				c.Set("k"+strconv.Itoa(i%100), i)
			}
			i++
		}
	})
}

func BenchmarkMutexBalance(b *testing.B) {
	c := &MutexMap{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				c.Get("k" + strconv.Itoa(i%100))
			} else {
				c.Set("k"+strconv.Itoa(i%100), i)
			}
			i++
		}
	})
}

func BenchmarkRWMutexBalance(b *testing.B) {
	c := &RWMutexMap{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				c.Get("k" + strconv.Itoa(i%100))
			} else {
				c.Set("k"+strconv.Itoa(i%100), i)
			}
			i++
		}
	})
}

func BenchmarkSyncBalance(b *testing.B) {
	c := &SyncMap{}
	for i := 0; i < 100; i++ {
		c.Set("k"+strconv.Itoa(i), i)
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				c.Get("k" + strconv.Itoa(i%100))
			} else {
				c.Set("k"+strconv.Itoa(i%100), i)
			}
			i++
		}
	})
}

func BenchmarkSyncWriteHeavyNewKeys(b *testing.B) {
	c := &SyncMap{}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i%100000) // 100k разных ключей
			c.Set(key, i)
			i++
		}
	})
}

func BenchmarkMutexWriteHeavyNewKeys(b *testing.B) {
	c := &MutexMap{data: map[string]int{}}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i%100000) // 100k разных ключей
			c.Set(key, i)
			i++
		}
	})
}
