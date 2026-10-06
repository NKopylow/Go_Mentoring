package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Тема: context, middleware, graceful shutdown.
//
// Задача 18. HTTP GET /orders/{id} с правильным контекстом.
//
// Слои: handler → service → repository.
// Middleware: RequestID, Timeout, Logging.
// Graceful shutdown по SIGINT (server.Shutdown + WaitGroup аудита).
//
// Тест: при таймауте — 504; request id доходит до репозитория (мок).

type Order struct {
	ID    string
	Total int64
}

type Repository interface {
	// Get имитирует медленную БД: ждать через select с ctx.Done(), не через time.Sleep.
	Get(ctx context.Context, id string) (Order, error)
}

type Service struct {
	repo Repository
	// TODO: WaitGroup для аудита
}

// Get вызывает repo.Get, затем асинхронно пишет аудит (~300ms).
// Аудит НЕ должен отменяться вместе с запросом → context.WithoutCancel(ctx),
// но request id из ctx должен сохраниться.
func (s *Service) Get(ctx context.Context, id string) (Order, error) {
	panic("not implemented")
}

// WithRequestID кладёт id в ctx через типизированный неэкспортируемый ключ (не string!).
func WithRequestID(ctx context.Context, id string) context.Context {
	panic("not implemented")
}

// RequestIDFrom достаёт request id из ctx; если нет — "".
func RequestIDFrom(ctx context.Context) string { panic("not implemented") }

// RequestID middleware: берёт X-Request-ID из заголовка или генерирует новый, кладёт в ctx.
func RequestID(next http.Handler) http.Handler { panic("not implemented") }

// Timeout(d) оборачивает r.Context() в WithTimeout.
// При DeadlineExceeded handler отвечает 504.
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	panic("not implemented")
}

// Logging логирует метод, путь, статус, длительность, request id.
// Для статуса нужна обёртка над http.ResponseWriter.
func Logging(next http.Handler) http.Handler { panic("not implemented") }

func main() {
	fmt.Println("not implemented")
}
