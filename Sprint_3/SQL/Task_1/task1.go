package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	err := client.Set(ctx, "user:3", "Петр Иванов", time.Second*10).Err()
	value, err := client.Get(ctx, "user:3").Result()

	fmt.Println(value, err)
}
