package main

import (
	"context"
	"fmt"
	"time"
)

// Handler должен вызвать Service, Service — Repository.

// Repository выполняется 5 секунд, но если HTTP-запрос отменился,
// вся цепочка должна завершиться.

func Handler(ctx context.Context) error {
	return Service(ctx)
}

func Service(ctx context.Context) error {
	return Repository(ctx)
}

func Repository(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()
	err := Handler(ctx)
	fmt.Println(err)
}
