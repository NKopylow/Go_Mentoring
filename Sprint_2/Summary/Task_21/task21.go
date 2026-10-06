package main

import (
	"context"
	"fmt"
	"time"
)

// Итоговая задача. Тема: всё из спринта.
//
// Задача 21. Простой Job Scheduler (worker pool + очередь + cancel + метрики).
//
// Минимум:
//   - фиксированный пул воркеров (WithWorkers);
//   - Submit / Cancel / Status / Start / Stop;
//   - паника в Job.Run → Failed, воркер жив;
//   - Stop дожидается running, новые не берёт.
//
// Тесты (-race): N задач все Done; Cancel pending; Cancel running (ctx.Done);
// Stop не стартует pending; нет утечки горутин после Stop.

type JobID string

type Status int

const (
	Pending Status = iota
	Running
	Done
	Failed
	Cancelled
)

// Job — единица работы. Run должен уважать ctx (отмена/таймаут).
type Job interface {
	ID() string
	Run(ctx context.Context) error
}

// JobFunc — адаптер функции к Job (как http.HandlerFunc).
type JobFunc struct {
	Name string
	Fn   func(ctx context.Context) error
}

func (j JobFunc) ID() string                    { return j.Name }
func (j JobFunc) Run(ctx context.Context) error { return j.Fn(ctx) }

type Metrics struct {
	Submitted, Running, Done, Failed, Cancelled int64
}

type Option func(*Scheduler)

// WithWorkers задаёт число воркеров (по умолчанию runtime.NumCPU()).
func WithWorkers(n int) Option { panic("not implemented") }

type Scheduler struct {
	// TODO: очередь, воркеры, статусы, cancel-функции, метрики
}

// New создаёт планировщик (ещё не запущен).
func New(opts ...Option) *Scheduler { panic("not implemented") }

// Start запускает воркеров. Повторный Start — ошибка.
func (s *Scheduler) Start(ctx context.Context) error { panic("not implemented") }

// Stop перестаёт брать новые задачи, ждёт running (ограничено ctx), останавливает воркеров.
func (s *Scheduler) Stop(ctx context.Context) error { panic("not implemented") }

// Submit ставит job в очередь → Pending. Если ctx уже отменён — Cancelled.
// Возвращает JobID (можно = job.ID()).
func (s *Scheduler) Submit(ctx context.Context, job Job) (JobID, error) {
	panic("not implemented")
}

// Cancel: Pending → убрать из очереди (Cancelled);
// Running → отменить её ctx; уже завершённую — false.
func (s *Scheduler) Cancel(id JobID) bool { panic("not implemented") }

// Status возвращает текущий статус задачи; false, если id неизвестен.
func (s *Scheduler) Status(id JobID) (Status, bool) { panic("not implemented") }

// Metrics — снимок атомарных счётчиков.
func (s *Scheduler) Metrics() Metrics { panic("not implemented") }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := New(WithWorkers(4))
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
