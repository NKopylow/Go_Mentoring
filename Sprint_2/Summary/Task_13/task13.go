package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Тема: примитивы синхронизации (RWMutex, Once, Cond/каналы), атомики, тестирование.
//
// Задача 13. Потокобезопасный кэш с TTL и дедупликацией загрузок (singleflight).
//
// Реализуй Cache[K comparable, V any] со следующими методами:
//
//   New[K, V](ttl time.Duration, opts ...Option) *Cache[K, V]
//   (c *Cache) Get(key K) (V, bool)                  — только чтение, под RLock
//   (c *Cache) Set(key K, value V)
//   (c *Cache) Delete(key K)
//   (c *Cache) GetOrLoad(ctx context.Context, key K, loader func(ctx context.Context) (V, error)) (V, error)
//   (c *Cache) Stats() Stats                         — Hits, Misses, Loads (атомики)
//   (c *Cache) Close()                               — останавливает фоновую очистку, идемпотентен (sync.Once)
//
// Требования к GetOrLoad:
//   - если значение в кэше и не протухло — вернуть его (hit);
//   - если нет — вызвать loader, НО: при N одновременных GetOrLoad с одним ключом loader
//     должен выполниться РОВНО ОДИН раз, остальные ждут его результат (singleflight).
//     Подсказка: map[K]*call, где call содержит done chan struct{}, val, err; либо sync.Once в call;
//   - ожидающие должны уважать ctx: если их контекст отменён — вернуть ctx.Err(), не дожидаясь загрузки;
//   - ошибка loader'а НЕ кэшируется;
//   - loader вызывается БЕЗ удержания мьютекса (иначе весь кэш встанет на время загрузки).
//
// Фоновая очистка (janitor): горутина по тикеру удаляет протухшие записи. Интервал задаётся
// Option'ом (WithCleanupInterval), по умолчанию ttl/2. Остановка — через Close().
//
// Тесты (go.mod + task13_test.go, запускать с -race):
//   - Set/Get/Delete базовые, табличный тест;
//   - протухание по TTL (используй маленький ttl, например 20ms);
//   - singleflight: 100 горутин одновременно вызывают GetOrLoad(одного ключа), loader с time.Sleep(50ms)
//     и atomic-счётчиком вызовов — ожидаем ровно 1 вызов и 100 одинаковых результатов;
//   - отмена контекста ожидающего во время загрузки;
//   - Close() дважды не паникует; после Close фоновых горутин нет (goroutine leak check);
//   - бенчмарк Get при 90% hit-rate с b.RunParallel.
//
// Вопрос для рефлексии (напиши комментарием в конце файла): почему здесь RWMutex выгоднее Mutex,
// и в каком сценарии RWMutex окажется медленнее?

type Stats struct {
	Hits, Misses, Loads int64
}

type Option func(*config)

type config struct {
	cleanupInterval time.Duration
}

func WithCleanupInterval(d time.Duration) Option {
	return func(c *config) { c.cleanupInterval = d }
}

type Cache[K comparable, V any] struct {
	mu sync.RWMutex
	// TODO: поля
}

func New[K comparable, V any](ttl time.Duration, opts ...Option) *Cache[K, V] {
	panic("not implemented")
}

func (c *Cache[K, V]) Get(key K) (V, bool) { panic("not implemented") }
func (c *Cache[K, V]) Set(key K, value V)  { panic("not implemented") }
func (c *Cache[K, V]) Delete(key K)        { panic("not implemented") }
func (c *Cache[K, V]) Stats() Stats        { panic("not implemented") }
func (c *Cache[K, V]) Close()              { panic("not implemented") }

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

// Ответ на вопрос про RWMutex:
//
