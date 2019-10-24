package main

import (
	"os"

	fifth "github.com/h8gi/fifth/lib"
)

func main() {
	i := fifth.NewInterpreter(os.Stdin)
	i.Repl()
}
