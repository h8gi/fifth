package fifth

import (
	"fmt"
	"testing"
)

func TestIfElseThen(t *testing.T) {
	i := NewInterpreter()
	var tests = []struct {
		input string
		want  string
	}{
		{"0 hello .s", "[1 2]"},
		{"1 hello .s", "[2]"},
		{"0 foo .s", "[1 2 2]"},
		{"1 foo .s", "[2]"},
		{"0 hoge .s", "[1 2]"},
		{"1 hoge .s", "[3 4]"},
	}

	i.SetString(`
: hello if 1 then 2 ;
: foo if 0 hello then 1 hello ;
: hoge if 1 2 else 3 4 then ;
see hello
see foo
see hoge`)
	if err := i.Run(); err != nil {
		t.Error(err.Error())
	}

	for _, test := range tests {
		i.DS.Clear()
		i.SetString(test.input)
		fmt.Println(test.input)
		if err := i.Run(); err != nil {
			t.Error(err.Error())
		}
		// string comparison
		if i.DS.String() != test.want {
			t.Errorf("%q => %v", test.input, i.DS.data)
		}
	}
}
