package main

import "testing"

func TestSquare(t *testing.T) {
	t.Run("Test square 5", func(t *testing.T) {
		if square(5) != 25 {
			t.Error("Result for 5 in 2 level is incorrect")
		}
	})
	t.Run("Test square 6", func(t *testing.T) {
		if square(6) != 36 {
			t.Error("Result for 6 in 2 level is incorrect")
		}
	})
}

func BenchmarkSquare5(b *testing.B) {
	for b.Loop() {
		square(5)
	}
}

func BenchmarkSquare6(b *testing.B) {
	for b.Loop() {
		square(6)
	}
}
