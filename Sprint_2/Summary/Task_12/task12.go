package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Тема: context + горутины.
//
// Задача 12. Кто ответил первым — тот победил.
//
// Представь: у тебя 3 копии одного сервиса. Ты спрашиваешь все сразу
// и берёшь ответ самой быстрой. Остальные запросы отменяешь.
//
// Шаги:
//   1. Реализуй FirstSuccess (см. ниже).
//   2. Запусти main — должен напечататься ответ A или B (не C, он всегда падает).
//   3. (по желанию) Напиши тест: все упали → ошибка; ctx с коротким таймаутом → ctx.Err().

var ErrNoFetchers = errors.New("no fetchers provided")

// FirstSuccess запускает все fns в отдельных горутинах ОДНОВРЕМЕННО.
//
// Что должна делать:
//  1. Если fns пустой — вернуть "", ErrNoFetchers.
//  2. Создать дочерний контекст: ctx, cancel := context.WithCancel(ctx)
//     (не забудь defer cancel()).
//  3. Для каждой fn запустить горутину: вызвать fn(childCtx), результат положить в канал.
//  4. Читать из канала:
//     - первая успешная (err == nil) → cancel() остальных, вернуть результат;
//     - если пришла ошибка — запомнить её, ждать дальше;
//     - если ошибок столько же, сколько fns — вернуть ошибку (можно errors.Join);
//     - если родительский ctx отменён — вернуть ctx.Err().
//
// Подсказка: канал результатов сделай буферизованным (len(fns)),
// чтобы горутины не застряли на отправке после того, как ты уже вышел из функции.
func FirstSuccess(ctx context.Context, fns ...func(ctx context.Context) (string, error)) (string, error) {
	panic("not implemented")
}

// replica — готовая "заглушка сервиса". Не меняй.
// fail=true → всегда возвращает ошибку после случайной задержки.
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
		replica("C", true),
	)
	fmt.Println("result:", res)
	fmt.Println("error:", err)
	fmt.Println("took:", time.Since(start))
}
