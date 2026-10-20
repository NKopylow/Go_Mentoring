package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Тема: тестирование + context.
//
// Задача 17. Retry — повторить операцию несколько раз.
//
// Цель: написать простую функцию Retry и покрыть её тестами.
//
// Шаги:
//   1. Реализуй Retry (см. ниже).
//   2. Добавь go.mod + task17_test.go с табличными тестами:
//        - успех с первой попытки;
//        - успех со 2-й попытки (op первые разы возвращает ошибку);
//        - все попытки провалились → ErrMaxAttempts;
//        - ctx отменён во время паузы → ctx.Err().
//   3. Запусти: go test -v -race .

var ErrMaxAttempts = errors.New("max attempts exceeded")

// Retry вызывает op до attempts раз.
//
// Алгоритм:
//
//	for attempt := 1; attempt <= attempts; attempt++ {
//	  1. err := op(ctx)
//	  2. если err == nil → вернуть nil (успех)
//	  3. если это была последняя попытка → вернуть ErrMaxAttempts
//	  4. иначе подождать delay, НО с проверкой отмены:
//
//	       select {
//	       case <-time.After(delay):
//	           // идём на следующую попытку
//	       case <-ctx.Done():
//	           return ctx.Err()
//	       }
//	}
//
// Не используй time.Sleep — он не умеет отменяться по ctx.
func Retry(ctx context.Context, attempts int, delay time.Duration, op func(ctx context.Context) error) error {
	panic("not implemented")
}

func main() {
	ctx := context.Background()
	calls := 0

	err := Retry(ctx, 5, 100*time.Millisecond, func(ctx context.Context) error {
		calls++
		fmt.Println("attempt", calls)
		if calls < 3 {
			return errors.New("temporary failure")
		}
		return nil
	})

	fmt.Println("calls:", calls, "err:", err) // ожидаем: calls=3 err=<nil>
}
