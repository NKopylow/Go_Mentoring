package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Тема: WaitGroup, Once, семафор (канал), context.
//
// Задача 14. Свой errgroup (без сторонних пакетов).
//
// Тесты (-race): успех; первая ошибка отменяет ctx; SetLimit(2) не даёт >2 активных;
// TryGo возвращает false при занятых слотах.

type PanicError struct {
	Value any
	Stack []byte
}

func (e PanicError) Error() string { return fmt.Sprintf("panic: %v", e.Value) }

type Group struct {
	// TODO: поля (wg, once, err, cancel, limit-семафор)
}

// WithContext возвращает группу и дочерний ctx.
// Дочерний ctx отменяется при первой ошибке или при Wait().
func WithContext(ctx context.Context) (*Group, context.Context) {
	panic("not implemented")
}

// SetLimit задаёт максимум одновременно работающих Go/TryGo.
// n <= 0 — без лимита. Вызов после старта задач — паника.
func (g *Group) SetLimit(n int) { panic("not implemented") }

// Go запускает fn в горутине. При лимите — блокируется, пока не освободится слот.
// Первая ошибка из fn фиксируется один раз (sync.Once) и отменяет дочерний ctx.
func (g *Group) Go(fn func() error) { panic("not implemented") }

// TryGo как Go, но не блокируется: если слота нет — false и fn не запускается.
func (g *Group) TryGo(fn func() error) bool { panic("not implemented") }

// Wait ждёт все задачи и возвращает первую ошибку (не Join).
// Панику внутри fn перехватить recover'ом → PanicError.
func (g *Group) Wait() error { panic("not implemented") }

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
