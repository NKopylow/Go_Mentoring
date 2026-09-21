package main

import (
	"fmt"
	"sync"
	"time"
)

// Напиши реализацию семафоры, которое ограничивает количество горутин.

func main() {
	semaphore := make(chan struct{}, 3)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			semaphore <- struct{}{}

			defer func() {
				<-semaphore
			}()

			fmt.Println("started:", id)

			time.Sleep(2 * time.Second) // чтобы увидеть, что по 3 отрабатывает

			fmt.Println("finished:", id)
			fmt.Println("------------------------------------------------------------")
		}(i)
	}

	wg.Wait()
}
