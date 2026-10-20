package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Тема: WaitGroup + context + sync.Once.
//
// Задача 14. Группа задач с отменой при первой ошибке.
//
// Похоже на упрощённый errgroup: запускаешь несколько задач,
// если одна упала — остальные получают сигнал через ctx и тоже останавливаются.
//
// Шаги:
//   1. Реализуй WithContext, Go, Wait.
//   2. Запусти main: задача 3 падает → остальные должны напечатать "cancelled".
//   3. Проверь: go run -race .

type Group struct {
	wg     sync.WaitGroup
	once   sync.Once
	err    error
	cancel context.CancelFunc
}

// WithContext создаёт группу и дочерний контекст.
//
//	child, cancel := context.WithCancel(parent)
//
// cancel сохрани в Group — он понадобится, когда какая-то задача вернёт ошибку.
// Верни (&Group{cancel: cancel}, child).
func WithContext(ctx context.Context) (*Group, context.Context) {
	panic("not implemented")
}

// Go запускает fn в новой горутине.
//
// Что должна делать горутина:
//  1. defer g.wg.Done()  (не забудь g.wg.Add(1) перед запуском)
//  2. вызвать err := fn()
//  3. если err != nil — один раз (через g.once.Do) сохранить err и вызвать g.cancel()
func (g *Group) Go(fn func() error) {
	panic("not implemented")
}

// Wait ждёт все горутины (g.wg.Wait) и возвращает первую ошибку (или nil).
// После Wait можно вызвать cancel, чтобы освободить ресурсы контекста.
func (g *Group) Wait() error {
	panic("not implemented")
}

var errBoom = errors.New("boom")

func main() {
	g, ctx := WithContext(context.Background())

	for i := 1; i <= 5; i++ {
		id := i // важно скопировать, иначе все горутины увидят одно i
		g.Go(func() error {
			if id == 3 {
				return fmt.Errorf("task %d: %w", id, errBoom)
			}
			select {
			case <-time.After(time.Duration(id) * 300 * time.Millisecond):
				fmt.Println("task", id, "done")
				return nil
			case <-ctx.Done():
				fmt.Println("task", id, "cancelled:", ctx.Err())
				return ctx.Err()
			}
		})
	}

	fmt.Println("Wait:", g.Wait())
}
