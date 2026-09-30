package main

import (
	"fmt"
	"sync"
)

// Fan-In / Merge

func Merge(results ...chan int) chan int {
	var wg sync.WaitGroup
	result := make(chan int)
	for _, channel := range results {
		wg.Go(func() {
			for value := range channel {
				result <- value
			}
		})
	}
	go func() {
		wg.Wait()
		close(result)
	}()
	return result
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		defer close(ch1)

		ch1 <- 1
		ch1 <- 3
		ch1 <- 5
	}()

	go func() {
		defer close(ch2)

		ch2 <- 2
		ch2 <- 4
	}()

	go func() {
		defer close(ch3)

		ch3 <- 7
		ch3 <- 8
	}()

	result := Merge(ch1, ch2, ch3)

	for value := range result {
		fmt.Println(value)
	}
}
