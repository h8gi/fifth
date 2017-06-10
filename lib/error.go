package fifth

import "fmt"

func UndefinedError(text string) error {
	return fmt.Errorf("Undefined word %q", text)
}

var UnderFlowError = fmt.Errorf("Stack underflow")

var EOFError = fmt.Errorf("Unexpected EOF")
