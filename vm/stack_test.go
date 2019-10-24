package vm

import "testing"

func TestStackPush(t *testing.T) {
	stackSize := 5
	s := newStack(stackSize)
	for i := 0; i < stackSize+1; i++ {
		err := s.push(Cell(i))
		if err != nil {
			if err != ErrStackOverflow {
				t.Error("illegal error:", err)
			}
		}
	}
}

func TestStackPop(t *testing.T) {
	s := newStack(5)
	expected := Cell(8)
	s.push(expected)
	actual, err := s.pop()
	if err != nil {
		t.Error("stack pop")
	}
	if actual != expected {
		t.Errorf("stack pop value: expected %d, but actual %d", expected, actual)
	}

	_, err = s.pop()
	if err == nil || err != ErrStackUnderflow {
		t.Error("can't detect stackoverflow")
	}
}
