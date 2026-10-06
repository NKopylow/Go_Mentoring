package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Тема: SOLID, DRY, KISS.
//
// Задача 19. Рефакторинг «божественного» OrderService.
//
// 1. Комментарием: какие принципы нарушены и где (SRP, OCP, DIP, ISP, DRY).
// 2. Отрефактори:
//    - интерфейсы PaymentProcessor / Notifier / OrderRepository (маленькие);
//    - оплата — Strategy + выбор по имени (новый способ без правки сервиса);
//    - уведомления — Composite из нескольких Notifier;
//    - зависимости через NewOrderService(...); методы принимают context.Context;
//    - без log.Fatal в бизнес-логике; sentinel-ошибки (ErrInsufficientFunds и т.п.).
// 3. Тесты на моках: успех; нет денег; ошибка репозитория → уведомление не шлётся.

type Order struct {
	ID       string
	UserID   string
	Email    string
	Items    []string
	Amount   float64
	Payment  string // "card" | "paypal" | "cash"
	Promo    string
	IsVIP    bool
	CreateAt time.Time
}

type OrderService struct {
	balances map[string]float64
}

func NewOrderService() *OrderService {
	return &OrderService{balances: map[string]float64{"u1": 1000, "u2": 10}}
}

// PlaceOrder валидирует заказ, считает скидку, списывает оплату, сохраняет, шлёт уведомления.
// После рефакторинга: делегирует оплату/сохранение/уведомления интерфейсам.
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

	amount := o.Amount
	if o.IsVIP {
		amount = amount * 0.9
	}
	if o.Promo == "SALE10" {
		amount = amount * 0.9
	}
	if o.Promo == "SALE20" {
		amount = amount * 0.8
	}
	if len(o.Items) > 5 {
		amount = amount * 0.95
	}

	switch o.Payment {
	case "card":
		fmt.Println("[card] charging", amount)
		if s.balances[o.UserID] < amount {
			return errors.New("insufficient funds")
		}
		s.balances[o.UserID] -= amount
	case "paypal":
		fmt.Println("[paypal] redirect to paypal, charging", amount)
		if s.balances[o.UserID] < amount {
			return errors.New("insufficient funds")
		}
		s.balances[o.UserID] -= amount
	case "cash":
		fmt.Println("[cash] will be paid on delivery", amount)
	default:
		return errors.New("unknown payment " + o.Payment)
	}

	f, err := os.OpenFile("orders.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(f, "%s;%s;%.2f;%s\n", o.ID, o.UserID, amount, strings.Join(o.Items, ","))
	f.Close()

	fmt.Printf("[email] to %s: order %s placed, total %.2f\n", o.Email, o.ID, amount)
	if o.IsVIP {
		fmt.Printf("[push] VIP %s: thanks for order %s\n", o.UserID, o.ID)
	}
	if amount > 500 {
		fmt.Printf("[slack] #sales: big order %s from %s: %.2f\n", o.ID, o.UserID, amount)
	}

	fmt.Println(time.Now().Format(time.RFC3339), "order placed", o.ID)
	return nil
}

func main() {
	s := NewOrderService()
	err := s.PlaceOrder(Order{
		ID: "o1", UserID: "u1", Email: "a@b.c", Items: []string{"coffee", "bagel"},
		Amount: 700, Payment: "card", Promo: "SALE10", IsVIP: true,
	})
	fmt.Println("err:", err)
}

// Нарушенные принципы:
// Как добавить "crypto" после рефакторинга:
