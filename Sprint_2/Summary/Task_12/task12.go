package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Тема: context, каналы, утечки горутин.
//
// Задача 12. Первый успешный ответ побеждает.
//
// В main: 3 реплики (одна всегда падает). Выведи победителя и время.
// Напиши тесты: победа быстрого; все упали; таймаут ctx. Запусти с -race.

var ErrNoFetchers = errors.New("no fetchers provided")

// FirstSuccess запускает все fns параллельно с общим дочерним контекстом.
//
// Поведение:
//   - len(fns) == 0 → ErrNoFetchers
//   - первый успешный (T, nil) → отменить остальных через ctx, вернуть результат
//   - все вернули ошибку → errors.Join(...ошибки)
//   - родительский ctx отменён раньше → ctx.Err()
//   - после возврата горутины не должны блокироваться навсегда (буфер канала / select)
func FirstSuccess[T any](ctx context.Context, fns ...func(ctx context.Context) (T, error)) (T, error) {
	panic("not implemented")
}

// WithRetry возвращает обёртку над fn: до attempts попыток с паузой delay между ними.
//
// Поведение:
//   - при успехе — сразу вернуть результат
//   - пауза через select { time.After / ctx.Done }, не через time.Sleep
//   - при отмене ctx во время паузы — вернуть ctx.Err()
func WithRetry[T any](attempts int, delay time.Duration, fn func(ctx context.Context) (T, error)) func(ctx context.Context) (T, error) {
	panic("not implemented")
}

func replica(name string, fail bool) func(ctx context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		d := time.Duration(100+rand.Intn(900)) * time.Millisecond
		select {
		case <-time.After(d):
			if fail {
				return "", fmt.Errorf("%s: internal error", name)
			}
			return fmt.Sprintf("%s answered in %v", name, d), nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	res, err := FirstSuccess(ctx,
		replica("A", false),
		replica("B", false),
		WithRetry(3, 50*time.Millisecond, replica("C", true)),
	)
	fmt.Println(res, err, time.Since(start))
}
