package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Тема: тестирование (табличные тесты, моки, fake clock, t.Parallel, t.Cleanup, бенчмарки,
// testing/synctest), DIP из SOLID.
//
// Задача 17. Retry с экспоненциальным backoff и его полноценное тестирование.
//
// Реализуй:
//
//   type Clock interface {
//       Now() time.Time
//       After(d time.Duration) <-chan time.Time
//   }
//
//   type Retrier struct { ... }   // зависит от Clock через интерфейс, а не от пакета time напрямую
//
//   func NewRetrier(opts ...RetryOption) *Retrier
//   func (r *Retrier) Do(ctx context.Context, op func(ctx context.Context) error) error
//
// Поведение Do:
//   - вызывает op до MaxAttempts раз (по умолчанию 3);
//   - пауза между попытками: BaseDelay * 2^(attempt-1), но не больше MaxDelay; с опциональным
//     джиттером (Jitter(func() float64) — тоже инъекция, чтобы тест был детерминированным);
//   - если op вернула ошибку, обёрнутую в Permanent(err) — прекращает повторы НЕМЕДЛЕННО и
//     возвращает исходную ошибку (errors.Is/As должны работать через обёртку);
//   - пауза прерывается при отмене ctx -> вернуть errors.Join(ctx.Err(), lastErr);
//   - после исчерпания попыток -> вернуть ErrMaxAttempts, обёрнутую вместе с последней ошибкой;
//   - опциональный хук OnRetry(func(attempt int, err error, nextDelay time.Duration)) для логирования.
//
// Тесты (go.mod + task17_test.go), ОБЯЗАТЕЛЬНО:
//   1. fakeClock — реализация Clock, в которой After возвращает канал, а тест сам "продвигает время"
//      (fake.Advance(d)); Do должна выполняться без реальных задержек. Все тесты на ретраи должны
//      проходить за миллисекунды.
//   2. Табличный тест на расчёт задержки (nextDelay) — 6+ кейсов, включая срезание по MaxDelay.
//   3. Тест: успех со 2-й попытки; Permanent прерывает; MaxAttempts исчерпан; отмена ctx в паузе.
//   4. Мок операции с подсчётом вызовов и записью аргументов — проверь, что в op передаётся
//      тот же ctx (или производный от него).
//   5. t.Parallel() во всех независимых тестах, t.Cleanup для освобождения ресурсов fakeClock.
//   6. Тест с testing/synctest (Go 1.25+): реализация Clock на реальном time, но запущенная внутри
//      synctest.Test — убедись, что виртуальное время ждать не нужно.
//   7. go test -cover: добейся покрытия >= 90% для task17.go. Посмотри HTML-отчёт
//      (go test -coverprofile=c.out && go tool cover -html=c.out) и закрой непокрытые ветки.
//   8. Бенчмарк Do с мгновенно успешной операцией — allocs/op должны быть минимальными.
//
// Пример использования смотри в main.

var ErrMaxAttempts = errors.New("max attempts exceeded")

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type permanentError struct{ err error }

func (e permanentError) Error() string { return e.err.Error() }
func (e permanentError) Unwrap() error { return e.err }

// Permanent помечает ошибку как неповторяемую.
func Permanent(err error) error { return permanentError{err: err} }

type RetryOption func(*Retrier)

func WithMaxAttempts(n int) RetryOption                          { panic("not implemented") }
func WithBaseDelay(d time.Duration) RetryOption                  { panic("not implemented") }
func WithMaxDelay(d time.Duration) RetryOption                   { panic("not implemented") }
func WithClock(c Clock) RetryOption                              { panic("not implemented") }
func WithJitter(rnd func() float64) RetryOption                  { panic("not implemented") }
func WithOnRetry(fn func(int, error, time.Duration)) RetryOption { panic("not implemented") }

type Retrier struct {
	// TODO: поля
}

func NewRetrier(opts ...RetryOption) *Retrier { panic("not implemented") }

func (r *Retrier) Do(ctx context.Context, op func(ctx context.Context) error) error {
	panic("not implemented")
}

func main() {
	r := NewRetrier(
		WithMaxAttempts(5),
		WithBaseDelay(100*time.Millisecond),
		WithMaxDelay(time.Second),
		WithOnRetry(func(attempt int, err error, next time.Duration) {
			fmt.Printf("attempt %d failed: %v, retry in %v\n", attempt, err, next)
		}),
	)

	calls := 0
	err := r.Do(context.Background(), func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("temporary failure")
		}
		return nil
	})
	fmt.Println("calls:", calls, "err:", err)
}
