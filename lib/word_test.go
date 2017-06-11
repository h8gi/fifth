package fifth

import (
	"fmt"
	"testing"
)

func TestIf(t *testing.T) {
	i := NewInterpreter()
	var tests = []struct {
		input string
		want  string
	}{
		{"0 hello", "[1 2]"},
		{"1 hello", "[2]"},
	}

	i.SetString(": hello if 1 then 2 ;")
	fmt.Println(": hello if 1 then 2 ;")
	if err := i.Run(); err != nil {
		t.Error(err.Error())
	}

	for _, test := range tests {
		i.DS.Clear()
		i.SetString(test.input)
		if err := i.Run(); err != nil {
			t.Error(err.Error())
		}

		if i.DS.String() != test.want {
			t.Errorf("%q => %v", test.input, i.DS.data)
		}
	}
}
