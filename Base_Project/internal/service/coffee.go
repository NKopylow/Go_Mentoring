package service

import (
	"errors"

	"coffee/internal/domain"
)

var (
	ErrDrinkNotFound       = errors.New("напиток не найден")
	ErrInsufficientPayment = errors.New("недостаточно оплаты")
	ErrInsufficientStock   = errors.New("не хватает ингредиентов")
	ErrInvalidParams       = errors.New("некорректные параметры")
)

type Storage interface {
	GetStock(ingredient domain.Ingredient) int
	SetStock(ingredient domain.Ingredient, quantity int)
	AddStock(ingredient domain.Ingredient, quantity int)
	GetAllStock() map[domain.Ingredient]int

	GetDrink(name string) (domain.Drink, bool)

	AddSale(sale domain.Sale)
	GetStats() (int, int)
}

type CoffeeService struct {
	storage Storage // Storage, а не MemoryStorage,
	// потому что наш сервис для кофе не должен зависеть от конкретной реализации
	// (MemoryStorage - конкретная реализация с методами),
	// а только от абстракции (в Storage можно будет доложить интерфейс с расширенным поведением)
	drinks []domain.Drink
}

func NewCoffeeService(
	storage Storage,
	drinks []domain.Drink,
) *CoffeeService {
	return &CoffeeService{
		storage: storage,
		drinks:  drinks,
	}
}

func (cs *CoffeeService) AddStock(
	ingredientName string,
	quantity int,
) error {
	ingredient := domain.Ingredient(ingredientName)

	if !domain.IsValidIngredient(ingredient) {
		return ErrInvalidParams
	}

	if quantity < 1 {
		return ErrInvalidParams
	}

	cs.storage.AddStock(ingredient, quantity)

	return nil
}

func (cs *CoffeeService) SetStock(
	ingredientName string,
	quantity int,
) error {
	ingredient := domain.Ingredient(ingredientName) // по факту просто приводим тип к string

	if !domain.IsValidIngredient(ingredient) {
		return ErrInvalidParams
	}

	if quantity < 0 {
		return ErrInvalidParams
	}

	cs.storage.SetStock(ingredient, quantity)

	return nil
}

func (s *CoffeeService) GetStock() map[domain.Ingredient]int {
	return s.storage.GetAllStock()
}

func (s *CoffeeService) Brew(
	drinkName string,
	payment int,
) ([]string, error) {
	drink, ok := s.storage.GetDrink(drinkName)

	if !ok {
		return nil, ErrDrinkNotFound
	}

	if payment < drink.Price {
		return nil, ErrInsufficientPayment
	}

	for ingredient, required := range drink.Ingredients {
		if s.storage.GetStock(ingredient) < required {
			return nil, ErrInsufficientStock
		}
	}

	for ingredient, required := range drink.Ingredients {
		current := s.storage.GetStock(ingredient)

		s.storage.SetStock(
			ingredient,
			current-required,
		)
	}

	s.storage.AddSale(domain.Sale{ // логичнее считать стоимость напитка, потому что юзер может внести больше необходимой суммы
		Drink:  drink.Name,
		Amount: payment,
	})

	return drink.Steps, nil
}

func (s *CoffeeService) Stats() (int, int) {
	return s.storage.GetStats()
}

func (s *CoffeeService) Menu() []domain.Drink {
	return s.drinks
}
