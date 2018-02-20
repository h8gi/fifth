package fifth

import "fmt"

type Word struct {
	Name          string
	IsImmediate   bool
	IsPrimitive   bool
	IsCompileOnly bool
	PrimBody      PrimBody
	Body          []*Word
	pc            int // program counter
}

type PrimBody func() error

func (w *Word) String() string {
	s := ""
	if w.IsPrimitive {
		s += fmt.Sprintf("primitive word: %q", w.Name)
	} else {
		s += fmt.Sprintf(": %s\n  ", w.Name)
		for _, bw := range w.Body {
			s += bw.Name + " "
		}
		s += "\n;"
	}
	if w.IsImmediate {
		s += " immediate"
	}
	return s
}

func (w *Word) Compile(bw *Word) {
	w.Body = append(w.Body, bw)
}
