package stack

type Stack struct {
	arr []int
}

func (s *Stack) Push(value int) { // через (s *Stack) указывается, что метод относится к типу Stack
	s.arr = append(s.arr, value)
}

func (s *Stack) Pop() int {
	last := len(s.arr) - 1
	value := s.arr[last]
	s.arr = s.arr[:last]
	return value
}

func (s *Stack) Slice() int {
	last := len(s.arr) - 1
	value := s.arr[last]
	s.arr = s.arr[:last]
	return value
}

func New(items []int) *Stack {
	return &Stack{arr: items}
}
