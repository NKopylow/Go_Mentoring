package main

import (
	"context"
	"fmt"
	"time"
)

// Тема: каналы (pipeline), context.
//
// Задача 11. Конвейер (pipeline) с отменой.
//
// Реализуй набор универсальных стадий конвейера. Все стадии:
//   - НЕ блокируются навсегда при отмене ctx (каждая отправка/чтение — через select с ctx.Done());
//   - закрывают свой выходной канал при завершении (входной канал закрыт или ctx отменён);
//   - не утекают горутинами.
//
// 1. Generate(ctx, nums ...int) <-chan int          — отправляет числа в канал.
// 2. Map[T, R any](ctx, in <-chan T, fn func(T) R) <-chan R
// 3. Filter[T any](ctx, in <-chan T, pred func(T) bool) <-chan T
// 4. FanOut[T, R any](ctx, in <-chan T, workers int, fn func(T) R) <-chan R
//      — fn выполняется в `workers` горутинах, результаты сливаются в один канал
//        (переиспользуй идею Merge из Task_8). Порядок результатов НЕ важен.
// 5. Take[T any](ctx, in <-chan T, n int) <-chan T  — пропускает первые n и закрывает выход.
//      Важно: после Take продюсеры выше по цепочке должны завершиться, а не зависнуть
//      (подумай, кто и как должен отменить контекст).
//
// В main собери цепочку: Generate -> Filter(чётные) -> FanOut(4, медленный квадрат с time.Sleep)
// -> Take(5) и выведи результат. Затем проверь с помощью runtime.NumGoroutine()
// (после небольшого time.Sleep), что все горутины завершились.
//
// Дополнительно: напиши тест, который запускает цепочку с context.WithTimeout
// и проверяет, что функция возвращается не позднее таймаута + небольшой дельты.
//
// Проверь с -race.

func Generate(ctx context.Context, nums ...int) <-chan int {
	panic("not implemented")
}

func Map[T, R any](ctx context.Context, in <-chan T, fn func(T) R) <-chan R {
	panic("not implemented")
}

func Filter[T any](ctx context.Context, in <-chan T, pred func(T) bool) <-chan T {
	panic("not implemented")
}

func FanOut[T, R any](ctx context.Context, in <-chan T, workers int, fn func(T) R) <-chan R {
	panic("not implemented")
}

func Take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	panic("not implemented")
}

func slowSquare(n int) int {
	time.Sleep(100 * time.Millisecond)
	return n * n
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	nums := make([]int, 0, 100)
	for i := 1; i <= 100; i++ {
		nums = append(nums, i)
	}

	even := Filter(ctx, Generate(ctx, nums...), func(n int) bool { return n%2 == 0 })
	squares := FanOut(ctx, even, 4, slowSquare)

	for v := range Take(ctx, squares, 5) {
		fmt.Println(v)
	}
}
