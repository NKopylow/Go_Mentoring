package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Тема: SOLID / DRY / KISS.
//
// Задача 19. Разрежь «божественный» сервис на части.
//
// Сейчас OrderService делает ВСЁ сам: валидацию, оплату, запись в файл, email.
// Его нельзя нормально протестировать (пишет в файл, печатает в консоль).
//
// Шаги:
//   1. Прочитай PlaceOrder и комментарием внизу напиши, какие принципы нарушены
//      (хотя бы SRP — одна ответственность, DIP — зависимость от деталей).
//   2. Выдели 2 интерфейса:
//
//        type Payment interface {
//            Charge(userID string, amount float64) error
//        }
//        type Notifier interface {
//            Notify(email, message string) error
//        }
//
//   3. Сделай простые реализации: CardPayment (как сейчас card/paypal),
//      ConsoleNotifier (fmt.Println вместо "отправки").
//   4. OrderService должен принимать Payment и Notifier в конструкторе
//      и вызывать их, а не делать всё внутри себя.
//   5. Убери log.Fatal — возвращай ошибки.
//   6. (по желанию) Напиши тест с fake Payment, который возвращает ошибку —
//      и проверь, что Notifier при этом НЕ вызывается.
//
// Не усложняй: KISS важнее «идеальной архитектуры». Двух интерфейсов достаточно.

type Order struct {
	ID      string
	UserID  string
	Email   string
	Items   []string
	Amount  float64
	Payment string // "card" | "paypal" | "cash"
}

type OrderService struct {
	balances map[string]float64
}

func NewOrderService() *OrderService {
	return &OrderService{balances: map[string]float64{"u1": 1000, "u2": 10}}
}

// PlaceOrder сейчас нарушает SRP. После рефакторинга должна:
//  1. Проверить заказ (userID, items, amount) — вернуть error, не log.Fatal.
//  2. Списать деньги через Payment.Charge.
//  3. Сохранить заказ (пока можно оставить запись в файл или упростить до fmt.Println).
//  4. Отправить уведомление через Notifier.Notify.
func (s *OrderService) PlaceOrder(o Order) error {
	if o.UserID == "" {
		log.Fatal("user id is empty")
	}
	if len(o.Items) == 0 {
		return errors.New("no items")
	}
	if o.Amount <= 0 {
		return errors.New("bad amount")
	}

	switch o.Payment {
	case "card", "paypal":
		fmt.Println("["+o.Payment+"] charging", o.Amount)
		if s.balances[o.UserID] < o.Amount {
			return errors.New("insufficient funds")
		}
		s.balances[o.UserID] -= o.Amount
	case "cash":
		fmt.Println("[cash] pay on delivery", o.Amount)
	default:
		return errors.New("unknown payment " + o.Payment)
	}

	f, err := os.OpenFile("orders.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(f, "%s;%s;%.2f;%s\n", o.ID, o.UserID, o.Amount, strings.Join(o.Items, ","))
	f.Close()

	fmt.Printf("[email] to %s: order %s placed\n", o.Email, o.ID)
	fmt.Println(time.Now().Format(time.RFC3339), "order placed", o.ID)
	return nil
}

func main() {
	s := NewOrderService()
	err := s.PlaceOrder(Order{
		ID: "o1", UserID: "u1", Email: "a@b.c",
		Items: []string{"coffee"}, Amount: 300, Payment: "card",
	})
	fmt.Println("err:", err)
}

// Какие принципы нарушены и где:
//
