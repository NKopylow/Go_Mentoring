package main

import (
	"fmt"
	"sync"
	"time"
)

// Используй sync.WaitGroup для ожидания завершения N горутин,
// каждая из которых выполняет таймер на случайное время.

func goroutineFabric(n int, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(i) * time.Millisecond)
		}()
	}
}

func main() {
	var wg sync.WaitGroup

	goroutineFabric(100, &wg)

	start := time.Now()
	wg.Wait()
	end := time.Since(start)

	fmt.Println("duration: ", end)
}
