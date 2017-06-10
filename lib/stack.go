package fifth

type Stack struct {
	data []int64
}

func (s *Stack) Push(elm int64) {
	s.data = append(s.data, elm)
}

func (s *Stack) Pop() (int64, error) {
	if len(s.data) == 0 {
		return 0, UnderFlowError
	}

	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, nil
}

func (s *Stack) Clear() {
	s.data = []int64{}
}

func (s *Stack) MakeBinFunc(op func(int64, int64) int64) func() error {
	return func() error {
		n, err := s.Pop()
		if err != nil {
			return err
		}
		m, err := s.Pop()
		if err != nil {
			return err
		}
		s.Push(op(n, m))
		return err
	}
}
