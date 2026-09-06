package main

type Box[T any] struct {
	Value T
}

func Filter[T any](arr []T, f func(T) bool) []T {
	result := []T{}

	for _, value := range arr {
		if f(value) {
			result = append(result, value)
		}
	}
	return result
}

func main() {
	arr := []int{1, 2, 3, 4, 5, 30, 40, 50, 55, 70}

	lessThanTen := func(el int) bool {
		return el < 10
	}

	arr = Filter(arr, lessThanTen)

	// box := Box[int]{Value: 10}
	// box1 := Box[string]{Value: "Hello"}
}
