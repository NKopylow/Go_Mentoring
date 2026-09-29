package main

import (
	"context"
	"fmt"
	"time"
)

func fetch(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("success")
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// fetch(ctx)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	start := time.Now()
	err := fetch(ctx)
	size := time.Since(start)

	fmt.Println(size)
	fmt.Println(err)
}
