package vm

import "fmt"

type stack struct {
	data    []Cell
	pointer int
}

var (
	ErrStackOverflow  = fmt.Errorf("stack overflow")
	ErrStackUnderflow = fmt.Errorf("stack underflow")
)

func newStack(size int) *stack {
	return &stack{
		data:    make([]Cell, size),
		pointer: -1,
	}
}

func (s *stack) String() string {
	str := ""
	str += fmt.Sprintf("<%d>", s.pointer+1)
	for i := s.pointer; i >= 0; i-- {
		str += fmt.Sprintf(" %v", s.data[i])
	}
	return str
}

// return top of stack
func (s *stack) tos() Cell {
	return s.data[s.pointer]
}

func (s *stack) push(elem Cell) error {
	s.pointer++
	if s.pointer >= len(s.data) {
		return ErrStackOverflow
	}
	s.data[s.pointer] = elem
	return nil
}

func (s *stack) pop() (Cell, error) {
	if s.pointer < 0 {
		return 0, ErrStackUnderflow
	}
	elem := s.tos()
	s.pointer--
	return elem, nil
}
