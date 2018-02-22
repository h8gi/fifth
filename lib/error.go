package fifth

import "fmt"

func UndefinedError(text string) error {
	return fmt.Errorf("Undefined word: %s", text)
}

var UnderFlowError = fmt.Errorf("Stack underflow")

var EOFError = fmt.Errorf("Unexpected EOF")

var NoLastWordError = fmt.Errorf("No last word")

var QuitError = fmt.Errorf("Bye!")

var UnstructuredError = fmt.Errorf("unstructured control flow")

func CompileOnlyError(text string) error {
	return fmt.Errorf("compile only word: %s", text)
}

var TypeError = fmt.Errorf("type mismatch")
