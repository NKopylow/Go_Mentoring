package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

// Тема: утечки памяти и горутин.
//
// Задача 16. Найди баги и почини.
//
// Программа специально написана плохо: память и число горутин растут без остановки.
//
// Как смотреть:
//   1. go run .  — в консоли каждую секунду печатаются goroutines и heapInuse.
//   2. Открой в браузере: http://localhost:6060/debug/pprof/goroutine?debug=1
//
// Что искать (в коде есть все эти проблемы):
//   a) Subscribe добавляет канал, но никто никогда не удаляет подписчика.
//   b) Publish шлёт в канал в новой горутине — если никто не читает, горутина зависает навсегда.
//   c) big[:8] хранит ссылку на весь массив 1MB → GC не может освободить память.
//   d) time.After в цикле создаёт таймер на каждую итерацию (лучше time.NewTimer / тикер).
//
// Задача: почини код так, чтобы goroutines и heapInuse перестали расти.
// Рядом с каждым фиксом напиши короткий комментарий: что было не так.
// В конце файла — «до» и «после» (примерно какие цифры видел).

type Event struct {
	ID      int
	Payload []byte
}

type Bus struct {
	mu   sync.Mutex
	subs []chan Event
}

func (b *Bus) Subscribe() <-chan Event {
	ch := make(chan Event)
	b.mu.Lock()
	b.subs = append(b.subs, ch)
	b.mu.Unlock()
	return ch
}

func (b *Bus) Publish(e Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		go func(ch chan Event) { ch <- e }(ch) // утечка?
	}
}

var (
	bus     = &Bus{}
	seen    = map[int][]byte{}
	seenMu  sync.Mutex
	samples [][]byte
)

func processBatch(id int) {
	ch := bus.Subscribe()

	big := make([]byte, 1<<20) // 1 MB
	big[0] = byte(id)
	samples = append(samples, big[:8]) // держит весь big?

	timeout := time.After(50 * time.Millisecond)
	for {
		select {
		case e := <-ch:
			seenMu.Lock()
			seen[e.ID] = e.Payload
			seenMu.Unlock()
		case <-timeout:
			return // а канал ch? подписчик остаётся в Bus!
		case <-time.After(10 * time.Millisecond):
			bus.Publish(Event{ID: id, Payload: make([]byte, 4096)})
		}
	}
}

func worker(jobs <-chan int) {
	for id := range jobs {
		processBatch(id)
	}
}

func main() {
	go func() { _ = http.ListenAndServe("localhost:6060", nil) }()

	jobs := make(chan int)
	for range 4 {
		go worker(jobs)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		var ms runtime.MemStats
		for range ticker.C {
			runtime.ReadMemStats(&ms)
			fmt.Printf("goroutines=%d heapInuse=%dMB\n", runtime.NumGoroutine(), ms.HeapInuse>>20)
		}
	}()

	for i := 0; ; i++ {
		jobs <- i
	}
}

// До исправлений:  goroutines≈???  heap≈???MB
// После исправлений: goroutines≈???  heap≈???MB
