package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

// Тема: утечки памяти и горутин, pprof, GC.
//
// Задача 16. Найди и почини утечки.
//
// В этой программе спрятано НЕСКОЛЬКО утечек (памяти и горутин). Программа поднимает pprof
// на localhost:6060 и крутит нагрузку.
//
// 1. Запусти программу и понаблюдай:
//      go tool pprof http://localhost:6060/debug/pprof/heap        (top, -inuse_space, list ...)
//      go tool pprof http://localhost:6060/debug/pprof/goroutine   (top, traces)
//      curl http://localhost:6060/debug/pprof/goroutine?debug=2    (стеки всех горутин)
//    Плюс вывод runtime.NumGoroutine() и MemStats.HeapInuse печатается раз в секунду.
//
// 2. Для КАЖДОЙ найденной утечки опиши комментарием рядом с кодом:
//    - в чём проблема (кто держит ссылку / кто блокируется навсегда);
//    - как ты её нашёл (какой профиль/команда);
//    - исправление.
//
// Подсказки, что искать (но не только это):
//    - горутина, которая навсегда блокируется на отправке в канал, который никто не читает;
//    - подписчики, которые никогда не отписываются; срез, из которого удаляют элементы
//      так, что backing array не отпускает ссылки;
//    - time.After в цикле / неостановленный time.Ticker;
//    - под-срез маленькой части большого буфера, который держит весь буфер живым;
//    - map, который только растёт (Go не уменьшает количество бакетов map при удалении ключей).
//
// 3. После исправлений число горутин и HeapInuse должны стабилизироваться. Приложи вывод "до/после"
//    комментарием в конце файла.
//
// 4. (Дополнительно) Добавь тест с проверкой утечки горутин: NumGoroutine до и после вызова
//    processBatch, с ожиданием через polling (не одним time.Sleep).

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

// Наблюдения ДО исправлений:
//
// Наблюдения ПОСЛЕ исправлений:
//
