package main

import (
	"fmt"
	"sync"
)

// Реализуй дедлок, а затем исправь его

func main() {
	var wg sync.WaitGroup
	channel := make(chan int, 1)
	wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	channel <- 500
	// }() deadlock!!!
	go func() {
		defer wg.Done()
		select {
		case channel <- 500:
			fmt.Println("sent")
		default:
			fmt.Println("channel is busy")
		}
	}()

	wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	channel <- 501
	// }() deadlock!!!

	go func() {
		defer wg.Done()
		select {
		case channel <- 501:
			fmt.Println("sent")
		default:
			fmt.Println("channel is busy")
		}
	}()
	wg.Wait()
	close(channel)
	fmt.Print("End: ")
}
