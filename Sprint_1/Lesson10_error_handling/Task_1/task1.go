package main

import "errors"

// Верни ErrEmpty, если name == ""

var ErrEmpty = errors.New("empty")

func read(name string) error {
	if name == "" {
		return ErrEmpty
	}

	return nil
}
