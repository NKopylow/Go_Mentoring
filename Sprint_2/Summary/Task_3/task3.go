package main

import (
	"fmt"
	"sync"
)

func Merge(channels ...chan int) <-chan int {
	fmt.Println(channels)
	var wg sync.WaitGroup
	result := make(chan int, len(channels))
	for _, channel := range channels {
		wg.Go(func() {
			res := <-channel
			result <- res
		})
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)
	result := Merge(ch1, ch2, ch3)

	ch1 <- 10
	ch2 <- 100
	ch3 <- 1000

	close(ch1)
	close(ch2)
	close(ch3)

	for value := range result {
		fmt.Println("value: ", value)
	}
}
