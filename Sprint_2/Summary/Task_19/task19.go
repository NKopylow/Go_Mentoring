package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Тема: SOLID, DRY, KISS, тестирование через моки.
//
// Задача 19. Рефакторинг "божественного" сервиса.
//
// Ниже — работающий, но плохо спроектированный OrderService. Его невозможно протестировать без
// реальной записи в файл и "отправки" писем, а добавление нового способа оплаты или уведомления
// требует правки switch'ей внутри сервиса.
//
// 1. Перед рефакторингом перечисли комментарием, какие принципы нарушены и ГДЕ именно
//    (ожидается минимум: SRP, OCP, DIP, ISP; плюс дублирование — DRY).
//
// 2. Отрефактори, сохранив поведение (вывод в консоль/файл может отличаться форматом, но
//    набор действий — тот же):
//    - выдели интерфейсы PaymentProcessor, Notifier, OrderRepository, Logger (маленькие, ISP);
//    - способы оплаты — отдельные реализации (Strategy) + реестр/фабрика по имени, чтобы новый
//      способ добавлялся без изменения сервиса (OCP);
//    - уведомления — несколько Notifier'ов, объединяемых в один (Composite), сервис не знает,
//      сколько их и какие;
//    - скидки — отдельная цепочка правил (DiscountRule), легко расширяемая;
//    - все зависимости внедряются через конструктор NewOrderService(...) (DIP);
//    - сервис принимает context.Context первым аргументом в публичных методах;
//    - ошибки: sentinel/типизированные ошибки (ErrInsufficientFunds, ValidationError{Field}),
//      обёртка с %w, никакого log.Fatal внутри бизнес-логики.
//
// 3. KISS/DRY: не создавай слои ради слоёв — если абстракция имеет одну реализацию и не нужна для
//    тестов, объясни, почему оставил или убрал её.
//
// 4. Тесты (go.mod + task19_test.go): моки для всех интерфейсов (ручные, без библиотек),
//    табличные тесты PlaceOrder: успех; недостаточно средств; невалидный заказ; ошибка репозитория
//    -> уведомление НЕ отправляется; проверка порядка вызовов (оплата до сохранения).
//    Покрытие сервиса >= 90%.
//
// 5. В конце файла напиши, как теперь добавить способ оплаты "crypto" и SMS-уведомление — сколько
//    файлов/строк нужно изменить в сервисе (ожидаемый ответ: ноль).

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

func (s *OrderService) PlaceOrder(o Order) error {
	// валидация
	if o.UserID == "" {
		log.Fatal("user id is empty")
	}
	if len(o.Items) == 0 {
		return errors.New("no items")
	}
	if o.Amount <= 0 {
		return errors.New("bad amount")
	}

	// скидки
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

	// оплата
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

	// сохранение
	f, err := os.OpenFile("orders.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(f, "%s;%s;%.2f;%s\n", o.ID, o.UserID, amount, strings.Join(o.Items, ","))
	f.Close()

	// уведомления
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

// Нарушенные принципы (до рефакторинга):
//
// Как добавить "crypto" и SMS после рефакторинга:
//
