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
		pointer: 0,
	}
}

// return top of stack
func (s *stack) tos() Cell {
	return s.data[s.pointer]
}

func (s *stack) push(elem Cell) error {
	if s.pointer >= len(s.data) {
		return ErrStackOverflow
	}
	s.data[s.pointer] = elem
	s.pointer++
	return nil
}

func (s *stack) pop() (Cell, error) {
	if s.pointer <= 0 {
		return 0, ErrStackUnderflow
	}
	s.pointer--
	elem := s.tos()
	return elem, nil
}
