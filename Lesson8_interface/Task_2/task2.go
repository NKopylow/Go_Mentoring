package main

import (
	"errors"
	"fmt"
)

// Создай тип, реализующий интерфейс error. Оберни экземпляр этой ошибки через fmt.Errorf("%w", …) и проверь, найдёт ли её errors.Is.

type NotFoundErr struct {
	Key string
}

func (e NotFoundErr) Error() string {
	return fmt.Sprintf("not found key: %s", e.Key)
}

func main() {
	err := NotFoundErr{Key: "user42"}
	wrapped := fmt.Errorf("storage: %w", err)

	fmt.Println(wrapped) // распечатать сообщение
	fmt.Println(errors.Is(wrapped, err))
}
