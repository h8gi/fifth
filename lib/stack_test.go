package fifth

import "testing"

func TestStack(t *testing.T) {
	s := &Stack{}
	if _, err := s.Pop(); err == nil {
		t.Error(`Can't detect stack underflow`)
	}
}
