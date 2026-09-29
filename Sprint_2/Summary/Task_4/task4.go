package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	ids := []int{
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
	}

	sem := make(chan struct{}, 5)

	var wg sync.WaitGroup
	var active atomic.Int32

	for _, id := range ids {
		wg.Go(func() {
			sem <- struct{}{}

			current := active.Add(1)
			fmt.Printf("Job %d START | active = %d\n", id, current)

			time.Sleep(2 * time.Second)

			current = active.Add(-1)
			fmt.Printf("Job %d END   | active = %d\n", id, current)

			<-sem
		})
	}

	wg.Wait()
}
