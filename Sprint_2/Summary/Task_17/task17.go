package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Тема: тестирование (табличные тесты, fake clock, t.Parallel), DIP.
//
// Задача 17. Retry с backoff.
//
// Retrier зависит от Clock через интерфейс (не от time напрямую) — чтобы тесты шли без реальных пауз.
//
// Тесты (go.mod + task17_test.go):
//   - fakeClock с Advance: успех со 2-й попытки; Permanent; MaxAttempts; отмена ctx;
//   - табличный тест расчёта задержки (срез по MaxDelay);
//   - t.Parallel() где можно.

var ErrMaxAttempts = errors.New("max attempts exceeded")

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type permanentError struct{ err error }

func (e permanentError) Error() string { return e.err.Error() }
func (e permanentError) Unwrap() error { return e.err }

// Permanent помечает ошибку как неповторяемую: Do сразу возвращает её без ретраев.
func Permanent(err error) error { return permanentError{err: err} }

type RetryOption func(*Retrier)

// WithMaxAttempts — сколько раз вызывать op (по умолчанию 3).
func WithMaxAttempts(n int) RetryOption { panic("not implemented") }

// WithBaseDelay — базовая пауза: delay = BaseDelay * 2^(attempt-1).
func WithBaseDelay(d time.Duration) RetryOption { panic("not implemented") }

// WithMaxDelay — верхняя граница паузы.
func WithMaxDelay(d time.Duration) RetryOption { panic("not implemented") }

// WithClock подставляет часы (в тестах — fake).
func WithClock(c Clock) RetryOption { panic("not implemented") }

type Retrier struct {
	// TODO: maxAttempts, baseDelay, maxDelay, clock
}

// NewRetrier собирает Retrier из опций (дефолты: 3 попытки, BaseDelay=100ms, MaxDelay=1s, real clock).
func NewRetrier(opts ...RetryOption) *Retrier { panic("not implemented") }

// Do вызывает op до MaxAttempts раз.
//
// Поведение:
//   - успех → nil
//   - Permanent(err) → сразу err (без ретраев)
//   - между попытками пауза через clock.After; при отмене ctx → errors.Join(ctx.Err(), lastErr)
//   - попытки кончились → ErrMaxAttempts + lastErr
func (r *Retrier) Do(ctx context.Context, op func(ctx context.Context) error) error {
	panic("not implemented")
}

func main() {
	r := NewRetrier(
		WithMaxAttempts(5),
		WithBaseDelay(100*time.Millisecond),
		WithMaxDelay(time.Second),
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
