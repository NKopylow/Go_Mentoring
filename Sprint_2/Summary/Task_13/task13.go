package main

import (
	"fmt"
	"sync"
	"time"
)

// Тема: Mutex / RWMutex.
//
// Задача 13. Простой потокобезопасный кэш со сроком жизни (TTL).
//
// Цель: несколько горутин читают и пишут в map без гонок.
//
// Шаги:
//   1. Добавь в Cache поля: data (map), ttl, mu (RWMutex).
//   2. Реализуй New, Get, Set, Delete.
//   3. Запусти с флагом race: go run -race .
//   4. (по желанию) Напиши TestCache_TTL: Set → сразу Get ok → sleep > ttl → Get false.

// entry — одна запись в кэше: значение + время, когда оно протухнет.
type entry struct {
	value     string
	expiresAt time.Time
}

type Cache struct {
	mu sync.RWMutex
	// TODO: data map[string]entry
	// TODO: ttl  time.Duration
}

// New создаёт пустой кэш. ttl — сколько живёт каждая запись после Set.
func New(ttl time.Duration) *Cache {
	panic("not implemented")
}

// Set сохраняет value по ключу key.
// Время протухания = сейчас + ttl кэша.
// Писать под Lock (не RLock!).
func (c *Cache) Set(key, value string) {
	panic("not implemented")
}

// Get возвращает (value, true), если ключ есть и ещё не протух.
// Если ключа нет или он протух — ("", false).
// Читать под RLock.
//
// Подсказка: протух = time.Now().After(entry.expiresAt)
func (c *Cache) Get(key string) (string, bool) {
	panic("not implemented")
}

// Delete удаляет ключ (если есть). Писать под Lock.
func (c *Cache) Delete(key string) {
	panic("not implemented")
}

func main() {
	cache := New(200 * time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			cache.Set("user", "Alice")
			v, ok := cache.Get("user")
			fmt.Println(v, ok)
		})
	}
	wg.Wait()

	time.Sleep(300 * time.Millisecond)
	v, ok := cache.Get("user")
	fmt.Println("after TTL:", v, ok) // ожидаем: "" false
}
