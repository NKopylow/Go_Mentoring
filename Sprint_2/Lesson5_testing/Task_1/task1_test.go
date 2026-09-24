package main

import "testing"

func TestSquare(t *testing.T) {
	t.Run("Test fill slice", func(t *testing.T) {
		if square(5) != 25 {
			t.Error("Result for 5 in 2 level is incorrect")
		}
	})
}
