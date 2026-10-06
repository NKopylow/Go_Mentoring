package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Тема: context (values, deadlines, propagation, WithoutCancel, AfterFunc), graceful shutdown,
// middleware (Decorator).
//
// Задача 18. HTTP-сервис с правильной работой контекста на всех слоях.
//
// Построй маленький сервис GET /orders/{id} по слоям handler -> service -> repository и реализуй:
//
// 1. Middleware RequestID: берёт X-Request-ID из заголовка или генерирует новый, кладёт в контекст
//    через ТИПИЗИРОВАННЫЙ неэкспортируемый ключ и функции RequestIDFrom(ctx)/WithRequestID(ctx, id).
//    Запрещено использовать string как тип ключа. Объясни комментарием, почему.
//
// 2. Middleware Timeout(d): оборачивает r.Context() в WithTimeout. Handler должен вернуть 504,
//    если контекст истёк (отличай context.DeadlineExceeded от context.Canceled -> 499/клиент ушёл).
//
// 3. Middleware Logging: логирует метод, путь, статус, длительность и request id. Для статуса
//    нужна обёртка над http.ResponseWriter (запоминает WriteHeader).
//
// 4. Repository.Get(ctx, id): имитирует медленную БД — time.Sleep заменить на select с ctx.Done().
//    Service.Get(ctx, id) вызывает репозиторий и логирует с request id из ctx.
//
// 5. Аудит: после успешного ответа сервис должен асинхронно записать событие в "audit log"
//    (горутина, 300ms). Эта запись НЕ должна отменяться, когда клиент ушёл/запрос завершился —
//    используй context.WithoutCancel(ctx), но значения (request id) должны сохраниться.
//    Проверь руками: curl с --max-time 0.1 — аудит всё равно должен записаться.
//
// 6. context.AfterFunc: зарегистрируй в handler функцию, которая при отмене контекста запроса
//    увеличивает atomic-счётчик "cancelledRequests" (вывести в GET /metrics).
//
// 7. Graceful shutdown: главный контекст — signal.NotifyContext(SIGINT, SIGTERM). По сигналу:
//    server.Shutdown с таймаутом 5с, дождаться завершения фоновых аудит-горутин (WaitGroup),
//    и только потом выйти. Проверь: запусти долгий запрос (id=slow), нажми Ctrl+C — запрос
//    должен успеть завершиться, новые соединения не принимаются.
//
// 8. Тесты (go.mod + task18_test.go): httptest.NewRecorder для middleware цепочки; проверка, что
//    при таймауте возвращается 504; что request id пробрасывается в репозиторий (мок репозитория
//    записывает RequestIDFrom(ctx)); что аудит выполняется даже при отменённом контексте запроса.
//
// Вопросы (ответить комментарием в конце файла):
//   - почему нельзя хранить context в структуре;
//   - чем опасен context.WithValue для передачи обязательных параметров;
//   - что произойдёт, если забыть вызвать cancel у WithTimeout.

type Order struct {
	ID    string
	Total int64
}

type Repository interface {
	Get(ctx context.Context, id string) (Order, error)
}

type Service struct {
	repo Repository
	// TODO: WaitGroup для аудита, логгер и т.д.
}

func (s *Service) Get(ctx context.Context, id string) (Order, error) { panic("not implemented") }

func RequestIDFrom(ctx context.Context) string                     { panic("not implemented") }
func WithRequestID(ctx context.Context, id string) context.Context { panic("not implemented") }

func RequestID(next http.Handler) http.Handler                { panic("not implemented") }
func Timeout(d time.Duration) func(http.Handler) http.Handler { panic("not implemented") }
func Logging(next http.Handler) http.Handler                  { panic("not implemented") }

func main() {
	// TODO: собрать зависимости, роутер, middleware-цепочку, запустить сервер с graceful shutdown.
	fmt.Println("not implemented")
}

// Ответы на вопросы:
//
