package main

import (
	"context"
	"fmt"
	"time"
)

// Итоговая задача спринта (все темы): планировщик задач (Job Scheduler).
//
// Задача 21. Реализуй библиотеку scheduler, объединяющую весь спринт. Структура — несколько файлов
// в этой же директории (scheduler.go, options.go, queue.go, metrics.go, *_test.go), package main
// можно оставить только для демо в task21.go, либо вынести ядро в подпакет internal/scheduler —
// на твой выбор, но обоснуй его комментарием.
//
// Функциональность:
//
//  1. Job — интерфейс { ID() string; Run(ctx context.Context) error }. Плюс адаптер JobFunc
//     (как http.HandlerFunc).
//
//  2. Scheduler:
//       New(opts ...Option) *Scheduler
//       Submit(ctx, job Job, opts ...JobOption) (JobID, error)   // JobOption: WithPriority, WithTimeout,
//                                                                 // WithRetry(policy), WithDelay(d)
//       Cancel(id JobID) bool
//       Status(id JobID) (Status, bool)                          // Pending | Running | Done | Failed | Cancelled
//       Start(ctx) / Stop(ctx) error                             // graceful: дождаться running, не брать новые
//       Subscribe(func(Event)) (unsubscribe func())               // Observer: события смены статуса
//       Metrics() Metrics                                         // атомики: submitted, running, done, failed,
//                                                                 // cancelled, avg/percentile duration
//
//  3. Очередь с приоритетами (container/heap) + отложенные задачи (WithDelay) — отдельная горутина-таймер,
//     которая перекладывает задачи в очередь, когда наступает время (один time.Timer на ближайшую задачу,
//     а не по таймеру на каждую — объясни, почему).
//
//  4. Worker pool с фиксированным числом воркеров (WithWorkers) и опциональным rate limit
//     (WithRateLimit(n per second) — реализовать через time.Ticker/token bucket на канале).
//
//  5. Retry policy — Strategy: NoRetry, FixedRetry(n, d), ExponentialRetry(n, base, max) (переиспользуй Task_17;
//     Clock должен внедряться, чтобы тесты шли без реальных задержек).
//
//  6. Context: Submit принимает ctx — если он отменён до старта, задача получает Cancelled;
//     каждая задача запускается с производным контекстом (таймаут + возможность Cancel(id));
//     Stop(ctx) ограничен таймаутом ctx.
//
//  7. Паника внутри Job.Run не должна ронять воркер: recover -> Failed с PanicError.
//
//  8. Память: результаты/статусы завершённых задач не должны копиться вечно — WithResultTTL(d)
//     и фоновая очистка (см. Task_13). Проверь heap-профилем, что при 1 млн коротких задач
//     память стабилизируется; sync.Pool для объектов задач — по желанию, с бенчмарком до/после.
//
//  9. SOLID: Scheduler зависит только от интерфейсов (Queue, Clock, Logger, RetryPolicy);
//     каждая сущность в своём файле; добавление новой стратегии ретрая или новой очереди
//     (FIFO вместо приоритетной) не требует изменений Scheduler.
//
// 10. Тестирование (все с -race):
//     - unit: очередь (приоритеты, стабильность для равных приоритетов), retry policy (табличные),
//       rate limiter (N задач за секунду с fake clock или synctest);
//     - интеграционные: 1000 задач / 8 воркеров — все Done, метрики сходятся;
//       Cancel pending-задачи; Cancel running-задачи (задача видит ctx.Done());
//       Stop дожидается running и не стартует pending; паника -> Failed;
//       Subscribe получает события в правильном порядке для одной задачи;
//     - goroutine leak check после Stop;
//     - бенчмарк Submit+Run no-op задач с ReportAllocs; цель — минимум аллокаций на задачу;
//     - go test -cover >= 85%.
//
// 11. Демо в main: HTTP-эндпоинты POST /jobs, GET /jobs/{id}, GET /metrics, pprof на :6060,
//     graceful shutdown по SIGINT (см. Task_18).
//
// Критерии приёмки: go vet, go test -race ./..., отсутствие утечек горутин, читаемый код
// с короткими функциями, интерфейсы объявлены на стороне потребителя.

type JobID string

type Status int

const (
	Pending Status = iota
	Running
	Done
	Failed
	Cancelled
)

type Job interface {
	ID() string
	Run(ctx context.Context) error
}

type JobFunc struct {
	Name string
	Fn   func(ctx context.Context) error
}

func (j JobFunc) ID() string                    { return j.Name }
func (j JobFunc) Run(ctx context.Context) error { return j.Fn(ctx) }

type Event struct {
	JobID JobID
	From  Status
	To    Status
	Err   error
	At    time.Time
}

type Metrics struct {
	Submitted, Running, Done, Failed, Cancelled int64
	AvgDuration                                 time.Duration
}

type Option func(*Scheduler)
type JobOption func(*jobOptions)

type jobOptions struct {
	// TODO: priority, timeout, delay, retry policy
}

type Scheduler struct {
	// TODO: поля
}

func New(opts ...Option) *Scheduler                  { panic("not implemented") }
func (s *Scheduler) Start(ctx context.Context) error { panic("not implemented") }
func (s *Scheduler) Stop(ctx context.Context) error  { panic("not implemented") }
func (s *Scheduler) Submit(ctx context.Context, job Job, opts ...JobOption) (JobID, error) {
	panic("not implemented")
}
func (s *Scheduler) Cancel(id JobID) bool                          { panic("not implemented") }
func (s *Scheduler) Status(id JobID) (Status, bool)                { panic("not implemented") }
func (s *Scheduler) Subscribe(fn func(Event)) (unsubscribe func()) { panic("not implemented") }
func (s *Scheduler) Metrics() Metrics                              { panic("not implemented") }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := New()
	if err := s.Start(ctx); err != nil {
		fmt.Println(err)
		return
	}
	defer s.Stop(context.Background())

	id, err := s.Submit(ctx, JobFunc{Name: "hello", Fn: func(ctx context.Context) error {
		fmt.Println("hello from job")
		return nil
	}})
	fmt.Println(id, err)
}
