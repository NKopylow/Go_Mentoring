package main

import (
	"context"
	"fmt"
	"time"
)

// Тема: паттерны проектирования в Go (Functional Options, Decorator/Middleware, Observer,
// Builder vs Options, Strategy), в связке с каналами и sync.
//
// Задача 20. Универсальная шина событий (EventBus) с middleware.
//
// Часть A. Functional Options.
//   Реализуй NewEventBus(opts ...Option) *EventBus с опциями:
//     WithBufferSize(n)        — размер буфера подписчика (по умолчанию 16);
//     WithWorkers(n)           — число горутин доставки (по умолчанию 1 => порядок событий сохраняется);
//     WithDropPolicy(p)        — что делать, если буфер подписчика полон: Block | DropNewest | DropOldest;
//     WithErrorHandler(fn)     — колбэк на ошибку обработчика (по умолчанию — логирование).
//   Опции должны валидироваться (n <= 0 -> дефолт или паника с понятным текстом — выбери и обоснуй).
//
// Часть B. Observer на каналах.
//   type Handler[T any] func(ctx context.Context, event T) error
//
//   func Subscribe[T any](bus *EventBus, topic string, h Handler[T]) (unsubscribe func())
//   func Publish[T any](ctx context.Context, bus *EventBus, topic string, event T) error
//
//   Требования:
//     - несколько подписчиков на один topic; каждый получает копию события;
//     - unsubscribe безопасен при вызове дважды и во время доставки (нет паники "send on closed channel");
//     - тип события проверяется: Publish[int] в topic, где подписчик Handler[string] -> ошибка
//       ErrTypeMismatch (подсказка: хранить подписчиков как any и делать type assertion, либо
//       хранить reflect.Type topic'а при первой подписке);
//     - Close(ctx) — перестаёт принимать Publish (ErrClosed), дожидается доставки уже принятых событий
//       или истечения ctx; идемпотентен.
//
// Часть C. Decorator (middleware) для обработчиков.
//   type Middleware[T any] func(Handler[T]) Handler[T]
//   func Chain[T any](h Handler[T], mws ...Middleware[T]) Handler[T]  // порядок — как у http middleware из урока
//   Реализуй middleware: Recover (паника -> ошибка), Timeout(d), Retry(n) (переиспользуй идею из Task_12/17),
//   Metrics (атомарные счётчики успех/ошибка/длительность, доступные через bus.Metrics()).
//
// Часть D. Strategy для DropPolicy — отдельные типы, реализующие интерфейс
//   type dropStrategy interface { offer(ch chan any, ev any) (delivered bool) }
//   чтобы добавить новую политику без изменений в EventBus.
//
// Тесты (go.mod + task20_test.go, с -race):
//   - доставка всем подписчикам; порядок при WithWorkers(1);
//   - unsubscribe во время активной публикации из 10 горутин — нет паник и гонок;
//   - DropNewest/DropOldest при медленном подписчике (проверь, какие события потерялись);
//   - Chain: порядок вызова middleware (записывай в срез под мьютексом);
//   - Recover ловит панику, Timeout возвращает context.DeadlineExceeded;
//   - Close дожидается in-flight событий; Publish после Close -> ErrClosed;
//   - бенчмарк Publish с 1 и 10 подписчиками, с ReportAllocs.
//
// Вопрос (комментарием в конце файла): когда в Go стоит выбирать Builder вместо Functional Options,
// и почему для конфигов обычно предпочитают Options?

type DropPolicy int

const (
	Block DropPolicy = iota
	DropNewest
	DropOldest
)

type Option func(*EventBus)

func WithBufferSize(n int) Option                              { panic("not implemented") }
func WithWorkers(n int) Option                                 { panic("not implemented") }
func WithDropPolicy(p DropPolicy) Option                       { panic("not implemented") }
func WithErrorHandler(fn func(topic string, err error)) Option { panic("not implemented") }

type Handler[T any] func(ctx context.Context, event T) error
type Middleware[T any] func(Handler[T]) Handler[T]

func Chain[T any](h Handler[T], mws ...Middleware[T]) Handler[T] { panic("not implemented") }

func Recover[T any]() Middleware[T]                { panic("not implemented") }
func Timeout[T any](d time.Duration) Middleware[T] { panic("not implemented") }
func Retry[T any](attempts int) Middleware[T]      { panic("not implemented") }

type EventBus struct {
	// TODO: поля
}

func NewEventBus(opts ...Option) *EventBus          { panic("not implemented") }
func (b *EventBus) Close(ctx context.Context) error { panic("not implemented") }

func Subscribe[T any](bus *EventBus, topic string, h Handler[T]) (unsubscribe func()) {
	panic("not implemented")
}

func Publish[T any](ctx context.Context, bus *EventBus, topic string, event T) error {
	panic("not implemented")
}

type OrderCreated struct {
	ID    string
	Total int64
}

func main() {
	bus := NewEventBus(WithBufferSize(8), WithDropPolicy(DropOldest))

	handler := Chain(
		func(ctx context.Context, e OrderCreated) error {
			fmt.Println("send email for", e.ID)
			return nil
		},
		Recover[OrderCreated](),
		Timeout[OrderCreated](time.Second),
		Retry[OrderCreated](3),
	)

	unsub := Subscribe(bus, "order.created", handler)
	defer unsub()

	_ = Publish(context.Background(), bus, "order.created", OrderCreated{ID: "o1", Total: 100})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = bus.Close(ctx)
}

// Builder vs Functional Options:
//
