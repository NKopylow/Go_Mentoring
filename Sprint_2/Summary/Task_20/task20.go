package main

import (
	"context"
	"fmt"
	"time"
)

// Тема: паттерны (Functional Options, Observer, Decorator).
//
// Задача 20. EventBus на каналах.
//
// Тесты (-race): доставка всем подписчикам; unsubscribe безопасен; Publish после Close → ErrClosed;
// Chain вызывает middleware в правильном порядке; Recover ловит панику.

type Option func(*EventBus)

// WithBufferSize — размер буфера канала подписчика (по умолчанию 16).
func WithBufferSize(n int) Option { panic("not implemented") }

type Handler[T any] func(ctx context.Context, event T) error
type Middleware[T any] func(Handler[T]) Handler[T]

// Chain оборачивает h middleware'ами (порядок как у http: последний mw — внешний).
func Chain[T any](h Handler[T], mws ...Middleware[T]) Handler[T] {
	panic("not implemented")
}

// Recover ловит панику в handler и превращает её в error.
func Recover[T any]() Middleware[T] { panic("not implemented") }

// Timeout отменяет ctx handler'а через d; при истечении — context.DeadlineExceeded.
func Timeout[T any](d time.Duration) Middleware[T] { panic("not implemented") }

type EventBus struct {
	// TODO: подписчики по topic, mu, closed-флаг
}

// NewEventBus создаёт шину с опциями.
func NewEventBus(opts ...Option) *EventBus { panic("not implemented") }

// Close перестаёт принимать Publish (дальше — ErrClosed), дожидается in-flight или ctx.
// Идемпотентен.
func (b *EventBus) Close(ctx context.Context) error { panic("not implemented") }

// Subscribe регистрирует h на topic. Возвращает unsubscribe:
// безопасен при повторном вызове и во время Publish (без "send on closed channel").
func Subscribe[T any](bus *EventBus, topic string, h Handler[T]) (unsubscribe func()) {
	panic("not implemented")
}

// Publish доставляет event всем подписчикам topic. Несколько подписчиков — каждый получает событие.
// Если шина закрыта — ErrClosed.
func Publish[T any](ctx context.Context, bus *EventBus, topic string, event T) error {
	panic("not implemented")
}

type OrderCreated struct {
	ID    string
	Total int64
}

func main() {
	bus := NewEventBus(WithBufferSize(8))

	handler := Chain(
		func(ctx context.Context, e OrderCreated) error {
			fmt.Println("send email for", e.ID)
			return nil
		},
		Recover[OrderCreated](),
		Timeout[OrderCreated](time.Second),
	)

	unsub := Subscribe(bus, "order.created", handler)
	defer unsub()

	_ = Publish(context.Background(), bus, "order.created", OrderCreated{ID: "o1", Total: 100})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = bus.Close(ctx)
}
