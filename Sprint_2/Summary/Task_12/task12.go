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
// Задача 12. "Гонка" запросов: первый успешный результат побеждает.
//
// Есть несколько реплик сервиса (fetchers). Нужно опросить их параллельно и вернуть
// ПЕРВЫЙ успешный ответ, отменив остальные.
//
// Реализуй:
//
//   func FirstSuccess[T any](ctx context.Context, fns ...func(ctx context.Context) (T, error)) (T, error)
//
// Требования:
//   - все fns запускаются одновременно с дочерним контекстом;
//   - как только один вернул (T, nil) — остальные отменяются через контекст, функция возвращает результат;
//   - если ВСЕ вернули ошибку — вернуть errors.Join всех ошибок;
//   - если ctx отменён/истёк раньше — вернуть ctx.Err();
//   - горутины НЕ должны утекать: даже после возврата FirstSuccess все fns должны иметь возможность
//     завершиться и записать результат, не блокируясь навсегда (подумай о размере буфера канала
//     или об отдельном select с ctx.Done() при отправке);
//   - len(fns) == 0 -> вернуть ошибку ErrNoFetchers.
//
// 2. Реализуй WithRetry(attempts int, delay time.Duration, fn) — обёртку (декоратор) над fetcher'ом,
//    которая повторяет вызов до attempts раз с паузой delay, но прерывается при отмене ctx
//    (пауза — через select с time.After/ctx.Done(), а не time.Sleep).
//
// 3. В main: 3 реплики со случайной задержкой 100–1000 мс, одна из них всегда падает.
//    Выведи победителя и время. Проверь, что runtime.NumGoroutine() возвращается к исходному
//    спустя ~1 с после ответа.
//
// 4. Напиши тесты: победа самого быстрого; все упали; таймаут контекста; отсутствие утечек
//    (сравни runtime.NumGoroutine() до и после с небольшим ожиданием). Запусти с -race.

var ErrNoFetchers = errors.New("no fetchers provided")

func FirstSuccess[T any](ctx context.Context, fns ...func(ctx context.Context) (T, error)) (T, error) {
	panic("not implemented")
}

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
