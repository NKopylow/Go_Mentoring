package main

import (
	"fmt"
	"sync"
)

// Тема: паттерн Observer (подписка на события).
//
// Задача 20. Простая шина событий.
//
// Идея: кто-то публикует событие по имени (topic), все подписчики этого topic получают его.
//
// Шаги:
//   1. Реализуй Subscribe и Publish.
//   2. Запусти main — оба обработчика должны напечатать сообщение.
//   3. Проверь: go run -race .
//
// Не нужно делать generics, middleware, DropPolicy и т.п. — только map + Mutex.

type Handler func(msg string)

type EventBus struct {
	mu   sync.Mutex
	subs map[string][]Handler // topic → список обработчиков
}

// NewEventBus создаёт шину с пустым map подписчиков.
func NewEventBus() *EventBus {
	panic("not implemented")
}

// Subscribe добавляет handler в список подписчиков topic.
// Вызывать под Lock.
func (b *EventBus) Subscribe(topic string, h Handler) {
	panic("not implemented")
}

// Publish вызывает ВСЕ handler'ы topic по очереди, передавая им msg.
// Если подписчиков нет — просто ничего не делать (не ошибка).
//
// Важно: пока идёшь по списку handler'ов, держи RLock/Lock,
// либо сначала скопируй срез handler'ов под Lock, а вызывай уже без блокировки
// (чтобы handler не завис, если внутри тоже обратится к шине).
func (b *EventBus) Publish(topic, msg string) {
	panic("not implemented")
}

func main() {
	bus := NewEventBus()

	bus.Subscribe("order.created", func(msg string) {
		fmt.Println("email:", msg)
	})
	bus.Subscribe("order.created", func(msg string) {
		fmt.Println("analytics:", msg)
	})
	bus.Subscribe("user.login", func(msg string) {
		fmt.Println("security:", msg)
	})

	bus.Publish("order.created", "order #42")
	// ожидаем:
	//   email: order #42
	//   analytics: order #42
	// (security не вызывается — другой topic)
}
