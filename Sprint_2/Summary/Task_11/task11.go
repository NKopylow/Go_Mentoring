package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Тема: каналы + context.
//
// Задача 11. Простой конвейер (pipeline).
//
// Цель: сделать цепочку из 3 стадий, где каждая читает из канала и пишет в другой.
//
// Общие правила для Generate / Filter / Map:
//   - работа идёт в отдельной горутине;
//   - когда вход закончился (канал закрыт) ИЛИ ctx отменён — закрыть свой выходной канал;
//   - при отправке/чтении используй select с case <-ctx.Done(), чтобы не зависнуть.
//
// Пример select при отправке:
//
//	select {
//	case out <- value:
//	case <-ctx.Done():
//		return
//	}
//
// В main уже собрана цепочка: Generate → Filter(чётные) → Map(квадрат).
// Проверь: go run .   и   go run -race .

// Generate запускает горутину и:
//  1. По очереди отправляет каждое число из nums в канал.
//  2. После всех чисел — закрывает канал.
//  3. Если ctx отменён раньше — сразу закрывает канал и выходит.
func Generate(ctx context.Context, nums ...int) <-chan int {
	jobs := make(chan int)
	var wg sync.WaitGroup

	select {
	case <-ctx.Done():
		close(jobs)
		return jobs
	default:
		wg.Go(func() {
			for _, value := range nums {
				jobs <- value
			}
		})

		go func() {
			wg.Wait()
			close(jobs)
		}()
		return jobs
	}
}

// Filter читает числа из in.
// Если pred(n) == true — отправляет n в выходной канал.
// Если pred(n) == false — пропускает число.
// Когда in закрыт или ctx отменён — закрывает выход и выходит.
func Filter(ctx context.Context, in <-chan int, pred func(int) bool) <-chan int {
	panic("not implemented")
}

// Map читает числа из in, применяет fn к каждому и пишет результат в выходной канал.
// Когда in закрыт или ctx отменён — закрывает выход и выходит.
func Map(ctx context.Context, in <-chan int, fn func(int) int) <-chan int {
	panic("not implemented")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	even := Filter(ctx, Generate(ctx, nums...), func(n int) bool {
		return n%2 == 0
	})
	squares := Map(ctx, even, func(n int) int {
		return n * n
	})

	for v := range squares {
		fmt.Println(v) // ожидаем: 4, 16, 36, 64, 100
	}
}
