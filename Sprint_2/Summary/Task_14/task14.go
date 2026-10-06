package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Тема: WaitGroup, Mutex/Once, каналы как семафор, context.
//
// Задача 14. Собственный errgroup.
//
// Реализуй аналог golang.org/x/sync/errgroup БЕЗ использования сторонних пакетов:
//
//   type Group struct { ... }
//
//   func WithContext(ctx context.Context) (*Group, context.Context)
//       — возвращает группу и производный контекст, который отменяется при ПЕРВОЙ ошибке
//         или после Wait().
//   func (g *Group) SetLimit(n int)
//       — максимум n одновременно работающих функций (n <= 0 — без лимита).
//         Вызывать после старта задач — паника (как в оригинале).
//   func (g *Group) Go(fn func() error)
//       — запускает fn в горутине. Если лимит достигнут — Go БЛОКИРУЕТСЯ до освобождения слота.
//   func (g *Group) TryGo(fn func() error) bool
//       — как Go, но не блокируется: если слота нет — возвращает false.
//   func (g *Group) Wait() error
//       — ждёт все функции, возвращает ПЕРВУЮ ошибку (не Join!), остальные игнорирует.
//         Первая ошибка фиксируется ровно один раз (sync.Once).
//
// Паника внутри fn: перехватить через recover, превратить в ошибку PanicError{Value, Stack}
// и вернуть из Wait (в оригинале Go 1.23+ она ре-паникуется в Wait — можешь реализовать
// любой из вариантов, но осознанно и с комментарием почему).
//
// Нулевое значение Group должно быть пригодно к использованию (как у sync.WaitGroup).
//
// Тесты (go.mod + task14_test.go, с -race):
//   - все успешны -> nil;
//   - первая ошибка отменяет контекст: вторая задача должна завершиться по ctx.Done() раньше своего sleep;
//   - SetLimit(2): запускаем 10 задач с atomic-счётчиком активных, максимум одновременно не превышает 2;
//   - TryGo возвращает false при занятых слотах;
//   - паника превращается в ошибку / ре-паникуется;
//   - zero value Group работает.
//
// В main: 5 "запросов" с задержками, третий падает — покажи, что остальные остановились по контексту.

type PanicError struct {
	Value any
	Stack []byte
}

func (e PanicError) Error() string { return fmt.Sprintf("panic: %v", e.Value) }

type Group struct {
	// TODO: поля
}

func WithContext(ctx context.Context) (*Group, context.Context) { panic("not implemented") }
func (g *Group) SetLimit(n int)                                 { panic("not implemented") }
func (g *Group) Go(fn func() error)                             { panic("not implemented") }
func (g *Group) TryGo(fn func() error) bool                     { panic("not implemented") }
func (g *Group) Wait() error                                    { panic("not implemented") }

var errBoom = errors.New("boom")

func main() {
	g, ctx := WithContext(context.Background())
	g.SetLimit(3)

	for i := 1; i <= 5; i++ {
		g.Go(func() error {
			if i == 3 {
				return fmt.Errorf("task %d: %w", i, errBoom)
			}
			select {
			case <-time.After(time.Duration(i) * 200 * time.Millisecond):
				fmt.Println("task", i, "done")
				return nil
			case <-ctx.Done():
				fmt.Println("task", i, "cancelled:", ctx.Err())
				return ctx.Err()
			}
		})
	}

	fmt.Println("Wait:", g.Wait())
}
