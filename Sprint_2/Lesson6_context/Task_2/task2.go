package main

import (
	"context"
	"fmt"
	"time"
)

// Задача: измени worker, чтобы после отмены контекста он гарантированно мог прекратить работу,
// даже если сейчас заблокирован на отправке в result.

func worker(ctx context.Context, result chan<- int) {
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			result <- i
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := make(chan int)

	go worker(ctx, result)

	fmt.Println(<-result)

	cancel()

	fmt.Println("done")
}
