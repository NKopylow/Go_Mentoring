package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	value atomic.Int32
	mu    *sync.RWMutex
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Inc()
} // увеличивает счетчик на 1

func (c *Counter) Value() int32 {
	return c.value.Load()
} // безопасно получить текущее значение счётчика

func main() {
	var wg sync.WaitGroup
	counter := Counter{
		value: *new(atomic.Int32),
		mu:    &sync.RWMutex{},
	}

	for i := 0; i < 100; i++ {
		wg.Go(func() {
			for j := 0; j < 1000; j++ {
				counter.Inc()
			}
		})
	}

	wg.Wait()

	fmt.Println(counter.Value())
}
