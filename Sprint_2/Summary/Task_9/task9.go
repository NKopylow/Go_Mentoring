package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	value atomic.Int32
}

func (c *Counter) Value() int32 {
	// c.mu.RLock()
	res := c.value.Load()
	// c.mu.RUnlock()
	return res
}

func (c *Counter) Inc() {
	// c.mu.RLock()
	c.value.Add(1)
	// c.mu.RUnlock()
}

func main() {
	c := Counter{
		value: *new(atomic.Int32),
	}
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			c.Inc()
		})
	}

	wg.Wait()

	fmt.Println(c.Value())
}
