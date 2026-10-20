package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Тема: context в HTTP.
//
// Задача 18. Таймаут для HTTP-обработчика.
//
// Цель: если запрос обрабатывается слишком долго — вернуть 504 и отменить работу.
//
// Шаги:
//   1. Реализуй slowHandler — имитация медленной БД через select (не Sleep!).
//   2. Реализуй Timeout middleware.
//   3. Собери сервер в main и проверь:
//        curl http://localhost:8080/fast   → 200 OK
//        curl http://localhost:8080/slow   → 504 (потому что Timeout=1s, а slow ждёт 3s)
//
// Подсказка по middleware — как в уроке Lesson7_patterns/middleware:
//
//	func Timeout(d time.Duration) func(http.Handler) http.Handler {
//	    return func(next http.Handler) http.Handler {
//	        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	            // создать ctx с таймаутом из r.Context()
//	            // передать в next через r.WithContext(ctx)
//	        })
//	    }
//	}

// slowHandler имитирует долгую работу (wait).
// Если ctx отменён раньше — пишет 504 и выходит.
// Если дождался — пишет 200 и "done".
//
// Каркас:
//
//	select {
//	case <-time.After(wait):
//	    w.WriteHeader(200)
//	    fmt.Fprintln(w, "done")
//	case <-r.Context().Done():
//	    http.Error(w, "timeout", http.StatusGatewayTimeout) // 504
//	}
func slowHandler(wait time.Duration) http.HandlerFunc {
	panic("not implemented")
}

// Timeout оборачивает next: у каждого запроса свой ctx с дедлайном d.
// Не забудь defer cancel().
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	panic("not implemented")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/fast", slowHandler(100*time.Millisecond))
	mux.Handle("/slow", slowHandler(3*time.Second))

	// весь сервер ограничен таймаутом 1 секунда
	handler := Timeout(time.Second)(mux)

	fmt.Println("listening on :8080")
	_ = http.ListenAndServe(":8080", handler)
}
