package main

import (
	"fmt"

	fifth "github.com/h8gi/fifth/lib"
)

func main() {
	fmt.Println("こんにちは!")
	i := fifth.NewInterpreter()
	i.Repl()
}
