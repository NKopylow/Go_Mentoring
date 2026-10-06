package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

// Тема: утечки памяти и горутин, pprof.
//
// Задача 16. Найди и почини утечки.
//
// Программа крутит нагрузку и печатает goroutines/heapInuse раз в секунду.
// pprof: localhost:6060/debug/pprof/{heap,goroutine}
//
// Найди утечки (их несколько), рядом с кодом кратко опиши проблему и фикс.
// После правок числа должны стабилизироваться. «До/после» — комментарием в конце.
//
// Подсказки: блокировка на send в канал без читателя; подписчики без Unsubscribe;
// time.After в цикле; под-срез, держащий большой буфер; map, которая только растёт.

type Event struct {
	ID      int
	Payload []byte
}

type Bus struct {
	mu   sync.Mutex
	subs []chan Event
}

// Subscribe добавляет подписчика и возвращает канал событий.
func (b *Bus) Subscribe() <-chan Event {
	ch := make(chan Event)
	b.mu.Lock()
	b.subs = append(b.subs, ch)
	b.mu.Unlock()
	return ch
}

// Publish отправляет событие всем подписчикам.
func (b *Bus) Publish(e Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		go func(ch chan Event) { ch <- e }(ch)
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

	big := make([]byte, 1<<20)
	big[0] = byte(id)
	samples = append(samples, big[:8])

	timeout := time.After(50 * time.Millisecond)
	for {
		select {
		case e := <-ch:
			seenMu.Lock()
			seen[e.ID] = e.Payload
			seenMu.Unlock()
		case <-timeout:
			return
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

// До:
// После:
