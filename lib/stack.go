package fifth

import "fmt"

type Value interface{}

type Stack struct {
	data []Value
}

func (s *Stack) String() string {
	return fmt.Sprintf("%v", s.data)
}

func (s *Stack) Push(elm Value) {
	s.data = append(s.data, elm)
}

func (s *Stack) Pop() (Value, error) {
	if len(s.data) == 0 {
		return 0, UnderFlowError
	}

	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, nil
}

func (s *Stack) Clear() {
	s.data = []Value{}
}

func (s *Stack) MakeBinFunc(op func(int, int) int) func() error {
	return func() error {
		n, err := s.Pop()
		if err != nil {
			return err
		}
		ni, ok := n.(int)
		if !ok {
			return TypeError
		}

		m, err := s.Pop()
		if err != nil {
			return err
		}
		mi, ok := m.(int)
		if !ok {
			return TypeError
		}

		s.Push(op(ni, mi))
		return err
	}
}
