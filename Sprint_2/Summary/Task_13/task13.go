package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Тема: RWMutex, Once, атомики, singleflight.
//
// Задача 13. Потокобезопасный кэш с TTL и singleflight.
//
// Тесты (go.mod + task13_test.go, -race):
//   - Set/Get/Delete; протухание по TTL (~20ms);
//   - 100 горутин GetOrLoad одного ключа → loader вызван ровно 1 раз;
//   - Close() дважды не паникует.

type Stats struct {
	Hits, Misses, Loads int64
}

type Cache[K comparable, V any] struct {
	mu sync.RWMutex
	// TODO: data, ttl, inflight-загрузки, атомики статистики
}

// New создаёт кэш. Значения живут ttl с момента Set/успешного GetOrLoad.
func New[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	panic("not implemented")
}

// Get возвращает (value, true), если ключ есть и не протух; иначе (zero, false).
// Только чтение — под RLock.
func (c *Cache[K, V]) Get(key K) (V, bool) { panic("not implemented") }

// Set сохраняет value с текущим временем (TTL считается отсюда).
func (c *Cache[K, V]) Set(key K, value V) { panic("not implemented") }

// Delete удаляет ключ, если он есть.
func (c *Cache[K, V]) Delete(key K) { panic("not implemented") }

// Stats возвращает снимки Hits/Misses/Loads (атомики).
func (c *Cache[K, V]) Stats() Stats { panic("not implemented") }

// Close идемпотентно останавливает ресурсы кэша (sync.Once), если они есть.
func (c *Cache[K, V]) Close() { panic("not implemented") }

// GetOrLoad:
//   - hit (есть и не протух) → вернуть значение, Hits++
//   - miss → вызвать loader РОВНО один раз на ключ (остальные ждут тот же результат)
//   - ошибка loader НЕ кэшируется; loader вызывать БЕЗ удержания мьютекса
//   - если ctx ожидающего отменён — вернуть ctx.Err(), не дожидаясь загрузки
func (c *Cache[K, V]) GetOrLoad(ctx context.Context, key K, loader func(ctx context.Context) (V, error)) (V, error) {
	panic("not implemented")
}

func main() {
	cache := New[string, string](time.Second)
	defer cache.Close()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			v, err := cache.GetOrLoad(context.Background(), "user:1", func(ctx context.Context) (string, error) {
				fmt.Println("loading...") // должно напечататься один раз
				time.Sleep(100 * time.Millisecond)
				return "Alice", nil
			})
			fmt.Println(v, err)
		})
	}
	wg.Wait()
	fmt.Printf("%+v\n", cache.Stats())
}
